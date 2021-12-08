use std::env;

fn main() {
    let args: Vec<String> = env::args().collect();
    let filename = args[1].clone();

    let entries = read_input(&filename);
    let entry_values = entries.iter().map(|entry| {
        let entry_data = entry
            .split('|')
            .into_iter()
            .map(String::from)
            .collect::<Vec<String>>();
        let digits = parse_into_digits(&entry_data[0]);
        let mapping = find_association(&digits);
        let input_digits = parse_into_digits(&entry_data[1]);
        input_digits
            .into_iter()
            .map(|d| mapping.iter().position(|&r| *r == d).unwrap().to_string())
            .collect::<Vec<String>>()
            .join("")
            .parse::<i32>()
            .unwrap()
    });

    let concat = entry_values.clone().fold(String::new(), |acc, entry| acc + &entry.to_string());
    println!("Part 1 - The count of 1,4,7 or 8 in the entry values is: {}", count_chars(&concat, vec!['1', '4', '7', '8']));
    println!("Part 2 - The sum of these entry values is: {}", entry_values.sum::<i32>());
}

fn read_input(filename: &str) -> Vec<String> {
    std::fs::read_to_string(filename)
        .unwrap()
        .split('\n')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(String::from)
        .collect()
}

fn parse_into_digits(string: &str) -> Vec<Digit> {
    string
        .split(' ')
        .into_iter()
        .filter(|s| !s.is_empty())
        .map(|s| Digit {
            segments: s.chars().collect(),
        })
        .collect()
}

fn count_chars(string: &str, chars: Vec<char>) -> i32 {
    string.chars().fold(0, |occ, c| if chars.contains(&c) { occ + 1 } else { occ })
}

#[derive(Debug, Eq)]
struct Digit {
    segments: Vec<char>,
}

impl Digit {
    fn len(&self) -> usize {
        self.segments.len()
    }

    fn contains(&self, other: &Self) -> bool {
        other
            .segments
            .iter()
            .fold(true, |eq, c| eq && self.segments.contains(&c))
    }
}
impl PartialEq for Digit {
    fn eq(&self, other: &Self) -> bool {
        self.segments.len() == other.segments.len()
            && self
                .segments
                .iter()
                .fold(true, |eq, c| eq && other.segments.contains(&c))
    }
}

// Start with 4, 7 and 8 identified
// 6 digits + the segment of 4 = 9
// 6 digits + the segment of 1 = 0
// 6 digits + Not 9 nor 0 = 6
// 5 digits + all but 1 digit of 6 = 5
// 5 digits + all overlap on segments of 9 and not 5 = 3
// The remaining is 2
fn find_association(digits: &Vec<Digit>) -> Vec<&Digit> {
    let one = digits.iter().find(|d| d.len() == 2).unwrap();
    let four = digits.iter().find(|d| d.len() == 4).unwrap();
    let seven = digits.iter().find(|d| d.len() == 3).unwrap();
    let eight = digits.iter().find(|d| d.len() == 7).unwrap();
    let nine = digits
        .iter()
        .find(|d| d.len() == 6 && d.contains(&four))
        .unwrap();
    let zero = digits
        .iter()
        .find(|d| d.len() == 6 && *d != nine && d.contains(&one))
        .unwrap();
    let six = digits
        .iter()
        .find(|d| d.len() == 6 && *d != nine && *d != zero)
        .unwrap();
    let five = digits
        .iter()
        .find(|d| d.len() == 5 && six.contains(d))
        .unwrap();
    let three = digits
        .iter()
        .find(|d| d.len() == 5 && *d != five && nine.contains(d))
        .unwrap();
    let two = digits
        .iter()
        .find(|&d| {
            d != zero
                && d != one
                && d != three
                && d != four
                && d != five
                && d != six
                && d != seven
                && d != eight
                && d != nine
        })
        .unwrap();
    vec![zero, one, two, three, four, five, six, seven, eight, nine]
}
