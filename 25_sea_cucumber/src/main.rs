use std::collections::HashSet;
use std::env::args;
use std::fs::read_to_string;

fn main() {
    let args = args().collect::<Vec<String>>();
    let filename = args[1].clone();
    let (mut efsc, mut sfsc, xmax, ymax) = read_map(&filename);
    let mut i = 0;
    loop {
        // println!("");
        // println!("After {} iter", i);
        // print_map(xmax, ymax, &efsc, &sfsc);
        i += 1;
        let mut sc_moved = false;
        let mut new_efsc = HashSet::new();
        let mut new_sfsc = HashSet::new();
        for (x, y) in efsc.iter() {
            let target_x = (x + 1) % xmax;
            if efsc.contains(&(target_x, *y)) || sfsc.contains(&(target_x, *y)) {
                new_efsc.insert((*x, *y));
                continue;
            }
            sc_moved = true;
            new_efsc.insert((target_x, *y));
        }
        efsc = new_efsc;
        for (x, y) in sfsc.iter() {
            let target_y = (y + 1) % ymax;
            if efsc.contains(&(*x, target_y)) || sfsc.contains(&(*x, target_y)) {
                new_sfsc.insert((*x, *y));
                continue;
            }
            sc_moved = true;
            new_sfsc.insert((*x, target_y));
        }
        sfsc = new_sfsc;
        if !sc_moved {
            break;
        }
    }
    println!("Let's touch down on {}!", i);
}

fn print_map(
    xmax: usize,
    ymax: usize,
    esfc: &HashSet<(usize, usize)>,
    ssfc: &HashSet<(usize, usize)>,
) {
    for y in 0..ymax {
        for x in 0..xmax {
            if esfc.contains(&(x, y)) {
                print!(">");
                continue;
            }
            if ssfc.contains(&(x, y)) {
                print!("v");
                continue;
            }
            print!(".");
        }
        println!("");
    }
}
fn read_map(
    filename: &str,
) -> (
    HashSet<(usize, usize)>,
    HashSet<(usize, usize)>,
    usize,
    usize,
) {
    let mut east_side_sc = HashSet::new();
    let mut south_side_sc = HashSet::new();

    let lines = read_to_string(filename).unwrap();
    let mut y = 0;
    let mut x = 0;
    for l in lines.split('\n').filter(|s| !s.is_empty()) {
        x = 0;
        for c in l.chars() {
            match c {
                'v' => {
                    south_side_sc.insert((x, y));
                }
                '>' => {
                    east_side_sc.insert((x, y));
                }
                _ => {}
            }
            x += 1;
        }
        y += 1;
    }
    (east_side_sc, south_side_sc, x, y)
}
