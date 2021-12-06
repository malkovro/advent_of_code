use criterion::{black_box, criterion_group, criterion_main, Criterion};

use lanternfish::fishlife;

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("Lantern Fish Input", |b| {
        b.iter(|| {
            let fishes_at_t0 = fishlife::fishes_at_t0(black_box("input.txt"));
            let pop = fishlife::population(fishes_at_t0, 256);
            let head_count = pop.iter().sum::<u128>();
            assert!(head_count == 1732731810807)
        })
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
