import json
import matplotlib.pyplot as plt

with open("./output/result-76597-0.json", encoding="utf8") as file:
    results = json.loads(file.read())
    print(len(results))

successful = list(filter(lambda r: r["Success"], results))
print(len(successful))


counts = {
    "{:.1f}".format(i / 10): 0 for i in range(0, 20 * 10 + 1)
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
