# Realna miara wpięcia: komenda liczona tylko, gdy jest wołana w pliku
# OSIĄGALNYM z klient/src/aplikacja.ts, z pominięciem pliku definicji kontraktu.
import json, pathlib, re, sys
korzen = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path('.')
b = korzen / 'budowa'
wszystkie = {x['typ'] for x in json.loads((b/'shared/contract.json').read_text(encoding='utf-8'))['komendy']}
stale = dict(re.findall(r"^\s*(\w+):\s*'([\w.]+)',", (b/'shared/contract.ts').read_text(encoding='utf-8'), re.M))
src = b/'klient/src'
def importy(p):
    out = set()
    for m in re.finditer(r"from '(\.[^']+)'", p.read_text(encoding='utf-8', errors='replace')):
        cel = (p.parent/m.group(1)).resolve()
        for kand in (cel, pathlib.Path(str(cel)+'.ts'), cel/'index.ts'):
            if kand.exists(): out.add(kand); break
    return out
start = (src/'aplikacja.ts').resolve(); osiag = {start}; stos = [start]
while stos:
    for q in importy(stos.pop()):
        if q not in osiag: osiag.add(q); stos.append(q)
wolane = set()
for p in osiag:
    if pathlib.Path(p).name in ('contract.ts', 'contract.go'): continue
    for m in re.finditer(r'Command\.(\w+)', pathlib.Path(p).read_text(encoding='utf-8', errors='replace')):
        if m.group(1) in stale: wolane.add(stale[m.group(1)])
print(f'{len(wolane)}/{len(wszystkie)}')
