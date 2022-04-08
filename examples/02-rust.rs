#[no_mangle]
pub extern "C" fn add_one(x: i32) -> i32 {
    x + 1
}

// rustup target add wasm32-unknown-unknown
// rustc --target wasm32-unknown-unknown -O --crate-type=cdylib 02-rust.rs -o 02-rust.wasm
// wasm-gc 02-rust.wasm 02-rust-small.wasm