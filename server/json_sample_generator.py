import json, csv, random

klsmartin_data = []
images = {}

with open('data/kls_martin.csv', 'r', encoding='utf-8') as f:
    reader = csv.reader(f)
    next(reader)  
    for row in reader:
        code, eng, viet, alternative, brand = row
        klsmartin_data.append([code, eng, viet, alternative, brand])

with open('data/martin_images.csv', 'r', encoding='utf-8') as f:
    reader = csv.reader(f)
    next(reader)  
    for row in reader:
        code, image, *rest = row
        images[code] = image 

results = []
for i in range(1, 21):
    index = random.randint(0, 5000 - 1)  
    code = klsmartin_data[index][0]
    results.append({
        "No.": i,
        "model": code,
        "description": klsmartin_data[index][1],
        "quantity": 1,
        "image": images.get(code, None),  
    })

with open('sample_testing_data.json', 'w', encoding='utf-8') as f:
    json.dump(results, f, indent=4, ensure_ascii=False)