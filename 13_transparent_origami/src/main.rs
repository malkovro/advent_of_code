mod origami;

fn main() {
    let args : Vec<String> = std::env::args().collect();
    let filename = args[1].clone();
    let number_of_active_points = origami::count_active_points(&filename);
    println!("Number of Active Points after one fold: {}", number_of_active_points);
}
