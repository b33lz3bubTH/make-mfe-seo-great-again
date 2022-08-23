// https://api.nrix.in/product-api/ui-resolver/domain/nrix/products/<BaseCode>
// QYv0jEfmJ_Rdvkn49BaVM_base00 -> Eg
use serde::{Deserialize, Serialize};
use axum::{
    Json,
    response::{Html, IntoResponse},
    http::{StatusCode, Uri, header::{self, HeaderMap, HeaderName}},
};

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
    "#;
static HTML_SEGMENT_2: &'static str = r#"
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
</body>
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
pub async fn main() {
     // Build our application by creating our router.
    let app = axum::Router::new()
        .route("/products/:pid",
            axum::routing::get(get_product)
        );

    // Run our application as a hyper server on http://localhost:3000.
    axum::Server::bind(&"0.0.0.0:3001".parse().unwrap())
        .serve(app.into_make_service())
        .await
        .unwrap();
}

fn get_default() -> String {
    let html = HTML_SEGMENT_1.to_owned() + &HTML_SEGMENT_2.to_owned().to_string();
    return html.to_string();
}
fn get_seo_html(result: String) -> String {
    let html = HTML_SEGMENT_1.to_owned() + &result.to_owned() + &HTML_SEGMENT_2.to_owned().to_string();
    return html.to_string();
}

async fn get_product(axum::extract::Path(id):axum::extract::Path<String>) -> Html<&'static str> {
    // Headers(vec![
    //     (HeaderName::from_static("Content-Type"), HeaderValue::from_static("text/html; charset=utf-8")),
    // ]);

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
                    // return Html(&get_default());
                },
            };
        }
        reqwest::StatusCode::UNAUTHORIZED => {
            println!("Need to grab a new token");
            // return Html(&get_default());
        }
        other => {
            // panic!("Uh oh! Something unexpected happened: {:?}", other);
            // return Html(&get_default());
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

    println!("result: {:?}", result);

    let text_html = get_seo_html(result);
    return Html(text_html.to_string());
}

/*

GET https://shop.nrix.in/products/yHhsDXv_egj9VGw84PAtt_base00


SPA: <html>

    <div style="none" id="_SEO_SHIT">
        <p>Product Name: <span> </span></p>
        <p>Product Description: <span> </span></p>
        <p>Product Price: <span> </span></p>
        <p>Product ID: <span> </span></p>

    <div>

    .....
    .....
    ....
    
    <script src="main.js" />
</html>


*/



/*
<script>
        window.onload = () => {
            setTimeout(()=>{
                document.querySelector("_SEO_SHIT").style.display = "none";
            }, 2000)
        }
        
    </script>
*/