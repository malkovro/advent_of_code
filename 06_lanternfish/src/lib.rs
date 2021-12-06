pub mod fishlife;

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_case_part_1_passes() {
        let fishes_at_t0 = fishlife::fishes_at_t0("test_input.txt");
        let population = fishlife::population(fishes_at_t0, 80);
        assert_eq!(population.iter().sum::<u128>(), 5934);
    }

    #[test]
    fn test_case_part_2_passes() {
        let fishes_at_t0 = fishlife::fishes_at_t0("test_input.txt");
        let population = fishlife::population(fishes_at_t0, 256);
        assert_eq!(population.iter().sum::<u128>(), 26984457539);
    }
}
