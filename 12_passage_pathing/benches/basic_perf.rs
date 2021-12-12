use criterion::{black_box, criterion_group, criterion_main, Criterion};

use passage_pathing::passage_pathing;

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("Passage Pathing Test Input", |b| {
        b.iter(|| passage_pathing::count_paths(black_box("input.txt")))
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
