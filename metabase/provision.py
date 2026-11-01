import json
import os
import re
import time
import urllib.error
import urllib.request


BASE_URL = os.getenv("METABASE_URL", "http://metabase:3000")
EMAIL = os.environ["METABASE_ADMIN_EMAIL"]
PASSWORD = os.environ["METABASE_ADMIN_PASSWORD"]
DB_PASSWORD = os.environ["METABASE_DB_PASSWORD"]
DB_NAME = os.getenv("POSTGRES_DB", "cmykarte")


class MetabaseAPIError(RuntimeError):
    def __init__(self, method, path, status):
        super().__init__(f"Metabase API request {method} {path} failed with HTTP {status}")
        self.status = status


def request(method, path, body=None, session=None):
    data = None if body is None else json.dumps(body).encode()
    headers = {"Content-Type": "application/json"}
    if session:
        headers["X-Metabase-Session"] = session
    req = urllib.request.Request(BASE_URL + path, data=data, headers=headers, method=method)
    try:
        with urllib.request.urlopen(req, timeout=30) as response:
            payload = response.read()
            return json.loads(payload) if payload else None
    except urllib.error.HTTPError as error:
        if path == "/api/setup" and error.code == 400:
            raise RuntimeError(
                "Metabase initial setup rejected the configured admin credentials. "
                "Set METABASE_ADMIN_EMAIL to a valid email address and "
                "METABASE_ADMIN_PASSWORD to a strong, uncommon password, then rerun "
                "the metabase-provisioner service."
            ) from None
        raise MetabaseAPIError(method, path, error.code) from None


def validate_setup_credentials():
    if not re.fullmatch(r"[^@\s]+@[^@\s]+\.[^@\s]+", EMAIL):
        raise ValueError("METABASE_ADMIN_EMAIL must be a valid email address.")
    if len(PASSWORD) < 12:
        raise ValueError(
            "METABASE_ADMIN_PASSWORD must be at least 12 characters and not a common password."
        )


def collection(payload):
    if isinstance(payload, dict):
        return payload.get("data", [])
    return payload


def wait_for_metabase():
    deadline = time.time() + 600
    while time.time() < deadline:
        try:
            return request("GET", "/api/session/properties")
        except (urllib.error.URLError, TimeoutError):
            time.sleep(5)
    raise TimeoutError("Metabase did not become available within 10 minutes")


def login(properties):
    if properties.get("setup-token"):
        validate_setup_credentials()
        try:
            request(
                "POST",
                "/api/setup",
                {
                    "token": properties["setup-token"],
                    "user": {
                        "first_name": "Metabase",
                        "last_name": "Admin",
                        "email": EMAIL,
                        "password": PASSWORD,
                    },
                    "prefs": {"site_name": "CMYKarte Analytics", "site_locale": "en"},
                },
            )
        except MetabaseAPIError as error:
            if error.status != 403:
                raise

    try:
        result = request("POST", "/api/session", {"username": EMAIL, "password": PASSWORD})
    except MetabaseAPIError as error:
        if error.status == 401:
            raise RuntimeError(
                "Could not log in to Metabase with METABASE_ADMIN_EMAIL and "
                "METABASE_ADMIN_PASSWORD. If Metabase was already set up, these must "
                "match its existing admin account; changing .env does not reset that account."
            ) from None
        raise
    return result["id"]


def get_or_create_database(session):
    databases = request("GET", "/api/database", session=session)
    for database in collection(databases):
        if database.get("name") == "CMYKarte analytics":
            return database["id"]
    database = request(
        "POST",
        "/api/database",
        {
            "name": "CMYKarte analytics",
            "engine": "postgres",
            "details": {
                "host": "postgres",
                "port": 5432,
                "dbname": DB_NAME,
                "user": "metabase_readonly",
                "password": DB_PASSWORD,
                "ssl": False,
            },
        },
        session,
    )
    return database["id"]


CARDS = [
    {
        "name": "Match outcomes and rates",
        "display": "bar",
        "query": """
            select
                status,
                matches,
                round(100.0 * matches / sum(matches) over (), 1) as match_percent
            from (
                select status, count(*) as matches
                from analytics_match_health
                group by status
            ) outcomes
            order by status
        """,
    },
    {
        "name": "Match end reasons",
        "display": "bar",
        "query": """
            select coalesce(end_reason, 'in_progress') as end_reason, count(*) as matches
            from analytics_match_health
            group by end_reason
            order by matches desc
        """,
    },
    {
        "name": "Match duration distribution",
        "display": "bar",
        "query": """
            select
                case
                    when duration_seconds < 60 then 'Under 1 min'
                    when duration_seconds < 180 then '1-3 min'
                    when duration_seconds < 300 then '3-5 min'
                    else '5+ min'
                end as duration_bucket,
                count(*) as matches
            from analytics_match_health
            where status = 'completed'
            group by duration_bucket
            order by min(duration_seconds)
        """,
    },
    {
        "name": "Win rate by side",
        "display": "bar",
        "query": """
            select
                winner_side,
                count(*) as wins,
                round(100.0 * count(*) / sum(count(*)) over (), 1) as win_percent
            from analytics_match_health
            where status = 'completed'
            group by winner_side
            order by winner_side
        """,
    },
    {
        "name": "Rounds per match distribution",
        "display": "bar",
        "query": """
            select rounds_played, count(*) as matches
            from analytics_match_health
            where status <> 'in_progress'
            group by rounds_played
            order by rounds_played
        """,
    },
    {
        "name": "First-player win rate",
        "display": "scalar",
        "query": """
            select round(
                100.0 * count(*) filter (where winner_id = first_player_id) / nullif(count(*), 0),
                1
            ) as first_player_win_percent
            from matches
            where status = 'completed'
        """,
    },
    {
        "name": "Card plays by side and cost",
        "display": "bar",
        "visualization_settings": {
            "graph.dimensions": ["card_cost", "player_side"],
            "graph.metrics": ["plays"],
            "graph.x_axis.scale": "ordinal",
        },
        "query": """
            select player_side, card_cost, count(*) as plays
            from analytics_card_usage
            group by player_side, card_cost
            order by card_cost, player_side
        """,
    },
    {
        "name": "Match wins by card cost",
        "display": "bar",
        "query": """
            select
                player_side,
                card_cost,
                count(*) as plays,
                round(100.0 * count(*) filter (where winner_side = player_side) / count(*), 1)
                    as winning_match_percent
            from analytics_card_usage
            where match_status = 'completed'
            group by player_side, card_cost
            order by card_cost, player_side
        """,
    },
    {
        "name": "Card outcomes by attributes",
        "display": "table",
        "query": """
            select
                card_type,
                card_cost,
                card_color,
                card_pattern,
                count(*) as plays,
                round(100.0 * count(*) filter (where winner_side = player_side) / count(*), 1)
                    as winning_match_percent
            from analytics_card_usage
            where match_status = 'completed'
            group by card_type, card_cost, card_color, card_pattern, player_side
            order by plays desc
            limit 50
        """,
    },
    {
        "name": "Most common consecutive card plays",
        "display": "table",
        "query": """
            with plays as (
                select
                    player_id,
                    player_side,
                    round,
                    card_type || ' cost ' || card_cost as card,
                    lag(card_type || ' cost ' || card_cost) over (
                        partition by match_id, player_id, round order by event_id
                    ) as previous_card
                from analytics_card_usage
            )
            select player_side, previous_card, card, count(*) as combinations
            from plays
            where previous_card is not null
            group by player_side, previous_card, card
            order by combinations desc
            limit 20
        """,
    },
    {
        "name": "Score progression",
        "display": "line",
        "query": """
            select match_id, event_id, score_rgb, score_cmy
            from (
                select
                    e.id as event_id,
                    e.match_id,
                    e.score_rgb,
                    e.score_cmy,
                    row_number() over (partition by e.match_id order by e.id) as event_number
                from match_events e
                where e.event_type in ('card_played', 'turn_ended')
            ) progression
            where event_number <= 100
            order by match_id, event_id
            limit 2000
        """,
    },
    {
        "name": "Turns and rounds per match",
        "display": "table",
        "query": """
            select
                m.id as match_id,
                m.status,
                m.rounds_played,
                count(e.id) as turns
            from matches m
            left join match_events e
                on e.match_id = m.id and e.event_type = 'turn_ended'
            group by m.id, m.status, m.rounds_played
            order by m.started_at desc
            limit 100
        """,
    },
    {
        "name": "Card plays per match",
        "display": "table",
        "query": """
            select match_id, count(*) as cards_played
            from analytics_card_usage
            group by match_id
            order by cards_played desc
            limit 100
        """,
    },
]

DASHBOARDS = {
    "Match Health": [
        "Match outcomes and rates",
        "Match end reasons",
        "Match duration distribution",
        "Rounds per match distribution",
        "Win rate by side",
        "First-player win rate",
    ],
    "Card Balance": [
        "Card plays by side and cost",
        "Match wins by card cost",
        "Card outcomes by attributes",
        "Most common consecutive card plays",
    ],
    "Gameplay Progression": [
        "Score progression",
        "Turns and rounds per match",
        "Card plays per match",
    ],
}


def get_or_create_cards(session, database_id):
    existing = collection(request("GET", "/api/card", session=session))
    cards = {card["name"]: card for card in existing if "name" in card}
    for definition in CARDS:
        if definition["name"] in cards:
            card = cards[definition["name"]]
            visualization_settings = definition.get("visualization_settings")
            if (
                visualization_settings is not None
                and card.get("visualization_settings") != visualization_settings
            ):
                updated = request(
                    "PUT",
                    f"/api/card/{card['id']}",
                    {"visualization_settings": visualization_settings},
                    session,
                )
                cards[definition["name"]] = updated
            continue
        card = request(
            "POST",
            "/api/card",
            {
                "name": definition["name"],
                "display": definition["display"],
                "dataset_query": {
                    "database": database_id,
                    "type": "native",
                    "native": {"query": definition["query"]},
                },
                "visualization_settings": definition.get("visualization_settings", {}),
            },
            session,
        )
        cards[definition["name"]] = card
    return cards


def get_or_create_dashboards(session, cards):
    existing = collection(request("GET", "/api/dashboard", session=session))
    dashboards = {dashboard["name"]: dashboard for dashboard in existing}
    for name, card_names in DASHBOARDS.items():
        dashboard = dashboards.get(name)
        if dashboard is None:
            dashboard = request(
                "POST",
                "/api/dashboard",
                {"name": name, "description": "Automatically provisioned CMYKarte analytics."},
                session,
            )
            dashboards[name] = dashboard

        current = request("GET", f"/api/dashboard/{dashboard['id']}", session=session)
        dashcards = []
        existing_card_ids = set()
        for item in current.get("dashcards", []):
            card_id = item.get("card_id") or (item.get("card") or {}).get("id")
            if card_id is None:
                continue
            existing_card_ids.add(card_id)
            dashcards.append(
                {
                    "id": item["id"],
                    "card_id": card_id,
                    "row": item["row"],
                    "col": item["col"],
                    "size_x": item["size_x"],
                    "size_y": item["size_y"],
                    "parameter_mappings": item.get("parameter_mappings", []),
                    "series": item.get("series", []),
                }
            )
        added_card = False
        for index, card_name in enumerate(card_names):
            card = cards[card_name]
            if card["id"] in existing_card_ids:
                continue
            dashcards.append(
                {
                    "id": -(index + 1),
                    "card_id": card["id"],
                    "row": (index // 2) * 6,
                    "col": (index % 2) * 6,
                    "size_x": 6,
                    "size_y": 6,
                    "parameter_mappings": [],
                    "series": [],
                }
            )
            added_card = True
        if added_card:
            request(
                "PUT",
                f"/api/dashboard/{dashboard['id']}",
                {"dashcards": dashcards},
                session,
            )


def main():
    properties = wait_for_metabase()
    session = login(properties)
    database_id = get_or_create_database(session)
    cards = get_or_create_cards(session, database_id)
    get_or_create_dashboards(session, cards)
    print("Metabase analytics database, cards, and dashboards are ready.")


if __name__ == "__main__":
    main()
