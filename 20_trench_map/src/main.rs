use std::collections::HashSet;
use std::env::args;

fn main() {
    let args: Vec<String> = args().collect();
    let filename = args[1].clone();
    let (algo, mut light_map) = read_input(&filename);
    // print_map(&light_map);

    let lighted_points_after_enhancements = (0..50)
        .map(|i| {
            let mut enhanced_light_map = HashSet::new();
            let ((x0, y0), (xm, ym)) = min_max(&light_map);
            for y in y0..=ym {
                for x in x0..=xm {
                    let p = (x, y);
                    let enhancement_key = get_enhancement_key(
                        p,
                        &light_map,
                        (x0, y0),
                        (xm, ym),
                        !algo[0] || i % 2 == 0,
                    );
                    if algo[enhancement_key] {
                        enhanced_light_map.insert(p);
                    }
                }
            }
            light_map = enhanced_light_map;
            light_map.len()
        })
        .collect::<Vec<usize>>();
    println!(
        "Lighted points after 2 enhacements: {}",
        lighted_points_after_enhancements[1]
    );
    println!(
        "Lighted points after 50 enhacements: {}",
        lighted_points_after_enhancements[49]
    );
}

// fn print_map(light_map: &HashSet<(i32, i32)>) {
//     let ((x0, y0), (xm, ym)) = min_max(light_map);
//     for y in y0..=ym {
//         for x in x0..=xm {
//             if light_map.contains(&(x, y)) {
//                 print!("#");
//             } else {
//                 print!(".");
//             }
//         }
//         println!("");
//     }
// }

fn get_enhancement_key(
    p: (i32, i32),
    light_map: &HashSet<(i32, i32)>,
    (x0, y0): (i32, i32),
    (xm, ym): (i32, i32),
    ocean_dark: bool,
) -> usize {
    let mut key = 0;
    for y in (p.1 - 1)..=(p.1 + 1) {
        for x in (p.0 - 1)..=(p.0 + 1) {
            key <<= 1;
            if (!ocean_dark && (x <= x0 || x >= xm || y <= y0 || y >= ym))
                || light_map.contains(&(x, y))
            {
                key += 1;
            }
        }
    }
    key
}

fn min_max(map: &HashSet<(i32, i32)>) -> ((i32, i32), (i32, i32)) {
    let x_min = map.iter().map(|p| p.0).min().unwrap() - 1;
    let y_min = map.iter().map(|p| p.1).min().unwrap() - 1;
    let x_max = map.iter().map(|p| p.0).max().unwrap() + 1;
    let y_max = map.iter().map(|p| p.1).max().unwrap() + 1;

    ((x_min, y_min), (x_max, y_max))
}

fn read_input(filename: &str) -> (Vec<bool>, HashSet<(i32, i32)>) {
    let file_string = std::fs::read_to_string(filename).unwrap();
    let mut lines = file_string.split('\n').into_iter();
    let enhancement_algo = lines.next().unwrap().chars().fold(vec![], |mut algo, c| {
        algo.push(match c {
            '.' => false,
            _ => true,
        });
        algo
    });

    let mut map = HashSet::new();
    let mut y = 0;
    while let Some(line) = lines.next() {
        if line.is_empty() {
            continue;
        }
        let mut x = 0;
        for c in line.chars() {
            if c == '#' {
                map.insert((x, y));
            }
            x += 1;
        }
        y += 1;
    }

    (enhancement_algo, map)
}
