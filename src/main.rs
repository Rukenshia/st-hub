extern crate actix_web;
use actix_web::{server, App, HttpRequest, Responder};

mod battle;

fn index(_req: &HttpRequest) -> &'static str {
    "Hello world!"
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
