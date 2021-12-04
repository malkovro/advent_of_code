use criterion::{black_box, criterion_group, criterion_main, Criterion};

use giantsquid::bingo;

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("Bingo Test Input", |b| {
        b.iter(|| bingo::play(black_box("input.txt"), black_box(true)))
    });
}

criterion_group!(benches, criterion_benchmark);
criterion_main!(benches);
