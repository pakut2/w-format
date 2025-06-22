const possibleMatches = new Map<number, string>([
    [3, "fizz"],
    [5, "buzz"],
]);

const fizzbuzz = (): void => {
    for (let i = 1; i <= 100; i++) {
        const matches = Array.from(possibleMatches.entries())
            .filter(([key]) => i % key === 0)
            .map(([, value]) => value);

        console.log(i, matches.join(" "));
    }
};

fizzbuzz();
