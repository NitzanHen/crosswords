import json
import matplotlib.pyplot as plt

results = []

for i in range(50):
    with open(f"./output/result-37967-{i}.json", encoding="utf8") as file:
        results.extend(json.loads(file.read())) 

print(len(results))

successful = list(filter(lambda r: r["Success"], results))
print(len(successful))


counts = {
    "{:.1f}".format(i / 10): 0 for i in range(0, 10 * 10 + 1)
}
for res in successful:
    bucket = "{:.1f}".format(res['Time'])
    counts[bucket] += 1


bars = list(sorted(counts.items(), key=lambda x: x[0]))
names = [bar[0] for bar in bars]
values = [bar[1] for bar in bars]


cumulative = [sum(values[:i]) for i in range(len(values))]

plt.bar(names, values)
plt.show()
