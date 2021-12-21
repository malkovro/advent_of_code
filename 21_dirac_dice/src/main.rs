use std::collections::HashMap;

fn main() {
    let player_1 = (1, 0);
    let player_2 = (10, 0);

    part_1(player_1, player_2);
    part_2(player_1, player_2);
}

fn part_1(mut player_1: (u8, u32), mut player_2: (u8, u32)) {
    let mut dice_roll_count: u32 = 0;
    let loosing_player;

    loop {
        let rolled_sum_value = rolled_sum(&mut dice_roll_count);
        let target_space = (1 + (player_1.0 as u32 + rolled_sum_value - 1) % 10) as u8;
        player_1 = (target_space, player_1.1 + target_space as u32);
        if player_1.1 >= 1000 {
            loosing_player = player_2;
            break;
        }
        let rolled_sum_value = rolled_sum(&mut dice_roll_count);
        let target_space = (1 + (player_2.0 as u32 + rolled_sum_value - 1) % 10) as u8;
        player_2 = (target_space, player_2.1 + target_space as u32);
        if player_2.1 >= 1000 {
            loosing_player = player_1;
            break;
        }
    }
    let score = loosing_player.1 as u32 * dice_roll_count * 3;
    println!("Loosing player * dice thrown = {}", score);
}

fn part_2(player_1: (u8, u32), player_2: (u8, u32)) {
    let mut universes: HashMap<((u8, u32), (u8, u32)), u128> =
        HashMap::from([((player_1, player_2), 1)]);
    let mut player_1_wins: u128 = 0;
    let mut keep_going = true;
    let mut player_2_wins: u128 = 0;
    while keep_going {
        keep_going = false;
        let mut next_round_universe = HashMap::new();
        for ((player_1, player_2), count) in universes.iter() {
            for roll_1_1 in 1..=3 {
                for roll_1_2 in 1..=3 {
                    for roll_1_3 in 1..=3 {
                        let target_space_1 =
                            1 + (player_1.0 + roll_1_1 + roll_1_2 + roll_1_3 - 1) % 10;
                        let player_1_next_score = player_1.1 + target_space_1 as u32;

                        if player_1_next_score >= 21 {
                            player_1_wins += count;
                            continue;
                        }
                        for roll_2_1 in 1..=3 {
                            for roll_2_2 in 1..=3 {
                                for roll_2_3 in 1..=3 {
                                    let target_space_2 =
                                        1 + (player_2.0 + roll_2_1 + roll_2_2 + roll_2_3 - 1) % 10;
                                    let player_2_next_score = player_2.1 + target_space_2 as u32;
                                    if player_2_next_score >= 21 {
                                        player_2_wins += count;
                                        continue;
                                    }
                                    keep_going = true;
                                    let next_round = next_round_universe
                                        .entry((
                                            (target_space_1, player_1_next_score),
                                            (target_space_2, player_2_next_score),
                                        ))
                                        .or_insert(0);
                                    *next_round += count;
                                }
                            }
                        }
                    }
                }
            }
        }
        universes = next_round_universe;
    }
    println!(
        "Player 1 wins : {} vs Player 2 wins: {}",
        player_1_wins, player_2_wins
    );
}

fn rolled_sum(dice_roll_count: &mut u32) -> u32 {
    let rolled_sum_value = ((*dice_roll_count * 3 + 3) + (*dice_roll_count * 3 + 1)) * 3 / 2;
    *dice_roll_count += 1;
    rolled_sum_value
}
