#[no_mangle]
pub extern "C" fn add(x: i32, y: i32) -> i32 {
    x + y
}

#[no_mangle]
fn sub(x: i32, y: i32) -> i32 {
    x - y
}

#[no_mangle]
pub extern "C" fn inc(x: i32) -> i32 {
    x + 1
}

#[no_mangle]
pub extern "C" fn show() {
    //
}

// rustup target add wasm32-unknown-unknown
// rustc --target wasm32-unknown-unknown -O --crate-type=cdylib 02-rust.rs -o 02-rust.wasm
// wasm-gc 02-rust.wasm 02-rust-small.wasm