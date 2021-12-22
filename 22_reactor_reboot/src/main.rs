use regex::Regex;
use std::env::args;

fn main() {
    let args: Vec<String> = args().collect();
    let filename = args[1].clone();
    let regions = read_input(&filename);

    let mut reboot_sequence_toggled_rods = 0;
    let mut init_sequence_toggled_rods = 0;
    let init_process_region = Region::On([
        CoordRange(-50, 50),
        CoordRange(-50, 50),
        CoordRange(-50, 50),
    ]);
    for i in 0..regions.len() {
        let region = regions[i];
        if let Region::Off(_) = region {
            continue;
        }

        if let Some(init_process_region_scoped) = region.truncate_region(&init_process_region) {
            init_sequence_toggled_rods +=
                init_process_region_scoped.volume_ignoring(&regions[i + 1..])
        }
        reboot_sequence_toggled_rods += region.volume_ignoring(&regions[i + 1..]);
    }
    println!("After init sequence: {}", init_sequence_toggled_rods);
    println!("After reboot sequence: {}", reboot_sequence_toggled_rods);
}

#[derive(Debug, PartialEq, Eq, Clone, Copy)]
struct CoordRange(i32, i32);

impl CoordRange {
    fn dist(&self) -> u128 {
        (self.1 - self.0 + 1).abs() as u128
    }

    fn sub_segment(&self, other: Self) -> Option<Self> {
        if self.1 < other.0 || self.0 > other.1 {
            None // other is not a sub segment of self
        } else {
            Some(CoordRange(
                std::cmp::max(other.0, self.0),
                std::cmp::min(self.1, other.1),
            ))
        }
    }
}

#[derive(Debug, PartialEq, Eq, Clone, Copy)]
enum Region {
    On([CoordRange; 3]),
    Off([CoordRange; 3]),
}

impl Region {
    fn ranges(&self) -> &[CoordRange; 3] {
        match self {
            Self::On(r) => r,
            Self::Off(r) => r,
        }
    }

    fn truncate_region(&self, other: &Self) -> Option<Self> {
        let sub_ranges = other.ranges();
        let ranges = self.ranges();

        // Only consider the sub segments on the 3 dimensions
        let truncated_coords = [
            ranges[0].sub_segment(sub_ranges[0])?,
            ranges[1].sub_segment(sub_ranges[1])?,
            ranges[2].sub_segment(sub_ranges[2])?,
        ];
        match self {
            Region::On(_) => Some(Region::On(truncated_coords)),
            Region::Off(_) => Some(Region::Off(truncated_coords)),
        }
    }

    fn volume_ignoring(&self, ignoring: &[Region]) -> u128 {
        self.volume() - self.overlap_volume(ignoring)
    }

    fn volume(&self) -> u128 {
        self.ranges().iter().fold(1, |acc, dim| dim.dist() * acc)
    }

    fn overlap_volume(&self, regions: &[Self]) -> u128 {
        let subregions = regions
            .iter()
            .filter_map(|region| self.truncate_region(region))
            .collect::<Vec<Region>>();
        (0..subregions.len())
            .map(|i| subregions[i].volume_ignoring(&subregions[i + 1..]))
            .sum()
    }
}

fn read_input(filename: &str) -> Vec<Region> {
    let re =
        Regex::new(r"(on|off)\sx=(-?\d+)..(-?\d+),y=(-?\d+)..(-?\d+),z=(-?\d+)..(-?\d+)").unwrap();
    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(|l| {
            let cap = re.captures_iter(l).next().unwrap();
            let coords = [
                CoordRange(cap[2].parse().unwrap(), cap[3].parse().unwrap()),
                CoordRange(cap[4].parse().unwrap(), cap[5].parse().unwrap()),
                CoordRange(cap[6].parse().unwrap(), cap[7].parse().unwrap()),
            ];
            match &cap[1] {
                "on" => Region::On(coords),
                _ => Region::Off(coords),
            }
        })
        .collect()
}
