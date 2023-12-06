names = []

with open("names.txt", 'r') as f:
    names = f.readlines()

folders = ['AS', 'AU', 'EU', 'LA']
n_names = int(len(names) / len(folders))

for i in range(len(folders)):
    with open(f"{folders[i]}/names.txt", 'w') as f:
        f.writelines(names[i*n_names:(i+1)*n_names])
