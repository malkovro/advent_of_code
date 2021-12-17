use std::env::args;
use std::iter::Iterator;

fn main() {
    let args: Vec<String> = args().collect();
    let x_min = args[1].parse::<i32>().unwrap();
    let x_max = args[2].parse::<i32>().unwrap();
    let y_min = args[3].parse::<i32>().unwrap();
    let y_max = args[4].parse::<i32>().unwrap();

    let zone = (x_min, x_max, y_min, y_max);
    let mut vx = x_max + 1;

    let mut highest_pos = 0;
    let mut count = 0;
    loop {
        vx -= 1;
        let mut vy = 500;
        if max_x_reach(vx) < x_min {
            break;
        }
        'y_iter: loop {
            vy -= 1;
            if vy < y_min {
                break 'y_iter;
            }
            if let Some(h) = is_valid((vx, vy), &zone) {
                count += 1;
                if highest_pos < h {
                    highest_pos = h;
                }
            }
        }
    }

    println!("Best Shot out of {} goes as high as {}", count, highest_pos);
}

fn max_x_reach(vx: i32) -> i32 {
    vx * (vx + 1) / 2
}

fn is_valid(velocity: (i32, i32), zone: &(i32, i32, i32, i32)) -> Option<i32> {
    let mut velocity = velocity;
    let mut pos = (0, 0);
    let mut highest_pos = 0;

    loop {
        let res = step(velocity, pos);
        velocity = res.0;
        pos = res.1;
        if landed(pos, zone) {
            return Some(highest_pos);
        }
        if overshot(pos, zone) {
            return None;
        }
        if highest_pos < pos.1 {
            highest_pos = pos.1;
        }
    }
}

fn landed(pos: (i32, i32), zone: &(i32, i32, i32, i32)) -> bool {
    pos.0 >= zone.0 && pos.0 <= zone.1 && pos.1 >= zone.2 && pos.1 <= zone.3
}

fn overshot(pos: (i32, i32), zone: &(i32, i32, i32, i32)) -> bool {
    pos.0 > zone.1 || pos.1 < zone.2
}

fn step(velocity: (i32, i32), position: (i32, i32)) -> ((i32, i32), (i32, i32)) {
    (
        (
            if velocity.0 > 0 { velocity.0 - 1 } else { 0 },
            velocity.1 - 1,
        ),
        (position.0 + velocity.0, position.1 + velocity.1),
    )
}
