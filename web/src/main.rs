use actix_files::Files;
use actix_web::dev::ServiceRequest;
use actix_web::{middleware, web, App, HttpServer};
use actix_web_openidconnect::ActixWebOpenId;
use crafting::dbal::get_database;
use crafting::migration::{Migrator, MigratorTrait};
use crafting_web::app::{App, DbConnection};
use dotenvy::dotenv;
use leptos::config::get_configuration;
use leptos::prelude::*;
use leptos_actix::{generate_route_list, LeptosRoutes};
use leptos_meta::MetaTags;
use log::LevelFilter;
use std::io;

#[cfg(feature = "ssr")]
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    env_logger::Builder::new()
        .filter_level(LevelFilter::Info)
        .parse_default_env()
        .init();

    let _ = dotenv();

    let conf = get_configuration(None).unwrap();
    let addr = conf.leptos_options.site_addr;
    // Generate the list of routes in your Leptos App
    let routes = generate_route_list(App);
    log::info!("listening on http://{addr}");

    let db = get_database().await.map_err(std::io::Error::other)?;

    Migrator::up(&db, None)
        .await
        .map_err(std::io::Error::other)?;

    let should_auth = |req: &ServiceRequest| {
        !req.path().starts_with("/pkg")
            && !req.path().starts_with("/assets")
            && req.method() != actix_web::http::Method::OPTIONS
    };

    let openid = ActixWebOpenId::builder(
        std::env::var("OIDC_CLIENT_ID").expect("OIDC_CLIENT_ID"),
        std::env::var("SERVER_HOST").expect("SERVER_HOST") + "auth_callback",
        std::env::var("OIDC_ISSUER").expect("OIDC_ISSUER"),
    )
    .allow_all_audiences(true)
    .redirect_on_error(true)
    .use_pkce(true)
    .should_auth(should_auth)
    .post_logout_redirect_url(std::env::var("SERVER_HOST").expect("SERVER_HOST"))
    .build_and_init()
    .await
    .map_err(io::Error::other)?;

    HttpServer::new(move || {
        let leptos_options = conf.leptos_options.clone();
        let site_root = &leptos_options.site_root;

        App::new()
            // serve JS/WASM/CSS from `pkg`
            .service(Files::new("/pkg", format!("{site_root}/pkg")))
            // serve other assets from the `assets` directory
            .service(Files::new("/assets", site_root.as_ref()))
            .wrap(openid.get_middleware())
            .configure(openid.configure_open_id())
            .leptos_routes(routes.to_owned(), {
                let leptos_options = leptos_options.clone();
                move || {
                    view! {
                        <!DOCTYPE html>
                        <html lang="de">
                            <head>
                                <meta charset="utf-8" />
                                <meta
                                    name="viewport"
                                    content="width=device-width, initial-scale=1"
                                />
                                <AutoReload options=leptos_options.clone() />
                                <HydrationScripts options=leptos_options.clone() />
                                <MetaTags />
                            </head>
                            <body>
                                <App />
                            </body>
                        </html>
                    }
                }
            })
            .app_data(web::Data::new(leptos_options.to_owned()))
            .app_data(DbConnection::new(db.clone()))
            .wrap(middleware::Compress::default())
    })
    .bind(&addr)?
    .run()
    .await
}
