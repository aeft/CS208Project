import random
import sympy

# Total number of numbers to generate
num_total = 5000
# Define the number of primes (1/5 of the total) and composites
num_prime = num_total // 5    # 200 primes
num_composite = num_total - num_prime  # 800 composites

# Overall allowed range for final numbers
overall_lower = 10**8    # 10^8
overall_upper = 10**16   # 10^16

# Probability for selecting the "small" branch
# (i.e., generating numbers that are not that large)
p_small = 0.3  # 30% chance to choose the small branch

primes_list = []
composites_list = []

# Generate 200 prime numbers
while len(primes_list) < num_prime:
    # With probability p_small, generate a smaller prime; otherwise, a larger one.
    if random.random() < p_small:
        lower = 10**8    # lower bound for small primes
        upper = 10**12   # upper bound for small primes
    else:
        lower = 10**12   # lower bound for large primes
        upper = 10**16   # upper bound for large primes
    try:
        # sympy.randprime returns a random prime in the interval [lower, upper)
        prime_num = sympy.randprime(lower, upper)
        primes_list.append(prime_num)
    except ValueError:
        # If no prime is found in the given interval, skip and try again.
        continue

# Generate 800 composite numbers (as semiprimes)
while len(composites_list) < num_composite:
    if random.random() < p_small:
        # Small branch: choose factors from a smaller range to get products between 10^8 and 10^12.
        factor_lower = 10**4    # 10^4
        factor_upper = 10**6    # 10^6; max product: 10^6 * 10^6 = 10^12
    else:
        # Large branch: choose factors from a larger range to allow products up to 10^16.
        factor_lower = 10**4    # 10^4
        factor_upper = 10**8    # 10^8; max product: 10^8 * 10^8 = 10^16
    # Generate two random prime factors
    p = sympy.randprime(factor_lower, factor_upper)
    q = sympy.randprime(factor_lower, factor_upper)
    composite = p * q
    # Check if the composite number falls within the overall range
    if overall_lower <= composite <= overall_upper:
        composites_list.append(composite)

# Merge the two lists and shuffle to mix primes and composites
numbers = primes_list + composites_list
random.shuffle(numbers)

# Option 1: Print each number on a new line
# for num in numbers:
#     print(num)

# Option 2: Write the numbers to a file (uncomment below to use)
with open("./api/client/test_numbers.txt", "w") as file:
    for num in numbers:
        file.write(f"{num}\n")
