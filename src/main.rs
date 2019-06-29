extern crate actix_web;
use actix_web::{server, App, HttpRequest, Responder, Json, Error};

mod battle;
use battle::BattleStatistics;

fn index(_req: &HttpRequest) -> Result<Json<BattleStatistics>, Error> {
    Ok(Json(BattleStatistics{damage: 1234, division: false, hits: 1, kills: 2, win: false}))
}

fn get_battles(_req: &HttpRequest) -> impl Responder {
    "[{'id': 0, 'statistics': { 'win': true, 'damage': 120000, 'kills': 2, 'hits': 287, 'divison': false }}]".replace("'", "\"")
}

fn main() {
    server::new(||
        App::new().resource("/", |r| r.f(index))
            .resource("/battles", |r| r.f(get_battles))
        )
        .bind("127.0.0.1:8088")
        .unwrap()
        .run();
}
