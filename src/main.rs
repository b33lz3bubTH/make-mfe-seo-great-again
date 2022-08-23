// https://api.nrix.in/product-api/ui-resolver/domain/nrix/products/<BaseCode>
// QYv0jEfmJ_Rdvkn49BaVM_base00 -> Eg
// ref : https://github.com/actix/examples/blob/master/https-tls/rustls/src/main.rs

use serde::{Deserialize, Serialize};
use actix_web::{get, web, App, HttpServer, Responder};
use actix_web::{http, middleware, HttpResponse ,http::header::ContentType};
use std::{fs::File, io::BufReader};
use actix_files::{Files};
use actix_web::dev::{ServiceResponse, ServiceRequest};
use actix_web_lab::web::redirect;

static HTML_SEGMENT_1: &'static str = r#"
        <!doctype html>
        <html lang="en">
        
        <head>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width,initial-scale=1">
        <title>homepage</title>
        <link rel="stylesheet" charset="UTF-8"
            href="https://cdnjs.cloudflare.com/ajax/libs/slick-carousel/1.6.0/slick.min.css" />
        <script src="https://ajax.googleapis.com/ajax/libs/webfont/1.6.26/webfont.js"></script>
        <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
        <script defer="defer" src="http://localhost:3001/main.js"></script>
        </head>
        
        <body>
            <div id="loading" style="overflow: hidden;">
            <div class="loading-container" style="text-align: center; width: 100vw; height: 100vh;"><img id="loading-image"
                src="https://api.nrix.in/media/c7ea07db-d00a-46e0-a18d-c3b61835b24d.png/thumb" alt="Loading..." style="width: 300px; height: 300px; 
                position: absolute; margin-top: -150px; 
                margin-left: -150px;
                top: 50%;
                left: 50%;" /></div>
        </div>
        <div id="app"></div>
        <script>$(window).on('load', function () {
            setTimeout(removeLoader, 2000);
          });
          function removeLoader() {
            $("\#loading").fadeOut(500, function () {
              $("\#loading").remove();
            });
          }
          // enable for deployment
          // console.log = () => {};</script>
        </body>
    "#;
static HTML_SEGMENT_2: &'static str = r#"
    
    </html>
"#;

#[derive(Serialize, Deserialize, Debug)]
struct APIResponse {
    message: String
}

#[derive(Serialize, Deserialize, Debug)] 
struct APIData {
    name: String,
    medias: Vec<String>,
    sellingPrice: u32,
    description: String
}


#[derive(Serialize, Deserialize, Debug)]
struct ProductData {
    api_response_info: APIResponse,
    data: APIData
}

#[tokio::main]
pub async fn main() -> std::io::Result<()> {
    HttpServer::new(|| {
        App::new().service(get_product).service(Files::new("/", "dist")
        .index_file("index.html")
        .default_handler(|req: ServiceRequest| {
            let (http_req, _payload) = req.into_parts();
            async {
                let response = actix_files::NamedFile::open("./dist/index.html")?.into_response(&http_req);
                Ok(ServiceResponse::new(http_req, response))
            }
        }))
    })
    .bind(("127.0.0.1", 3001))?
    .run()
    .await
   
}

fn get_default() -> String {
    let html = HTML_SEGMENT_1.to_owned() + &HTML_SEGMENT_2.to_owned().to_string();
    return html.to_string();
}
fn get_seo_html(result: String) -> String {
    let html = HTML_SEGMENT_1.to_owned() + &result.to_owned() + &HTML_SEGMENT_2.to_owned().to_string();
    return html.to_string();
}
#[get("/products/{id}")]
async fn get_product(id: web::Path<String>) -> HttpResponse {
    let response = reqwest::get("https://api.nrix.in/product-api/ui-resolver/domain/nrix/products/".to_string() + &id)
    .await
    .unwrap();

    let product_data: ProductData;
    match response.status() {
        reqwest::StatusCode::OK => {
            match response.json::<ProductData>().await {
                Ok(parsed_product) => {
                    product_data = parsed_product;
                },
                Err(_) => {
                    println!("Hm, the response didn't match the shape we expected.");
                    return HttpResponse::Ok().content_type(ContentType::html()).body(
                        get_default()
                    )
                },
            };
        }
        reqwest::StatusCode::UNAUTHORIZED => {
            println!("Need to grab a new token");
            return HttpResponse::Ok().content_type(ContentType::html()).body(
                get_default()
            )
        }
        other => {
            return HttpResponse::Ok().content_type(ContentType::html()).body(
                get_default()
            )

        }
    };

    let result = format!("
        <div id=\"_SEO_SHIT\">
            <p>Product Name: <span> {} </span></p>
            <p>Product Description: <span> {} </span></p>
            <p>Product Price: <span> {} </span></p>
            <p>Product ID: <span>{}</span></p>
        </div>
        ", 
        product_data.data.name,
        product_data.data.description,
        product_data.data.sellingPrice, &id);

    // println!("result: {:?}", result);

    let text_html = get_seo_html(result);
    return HttpResponse::Ok().content_type(ContentType::html()).body(
        text_html.to_string()
    )
    
}