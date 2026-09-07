# Pokrycie wiązań: komenda liczona tylko wtedy, gdy któreś wiązanie klienta
# faktycznie ją woła. Rejestr komend i wytwór kontraktu są pomijane — wymieniają
# wszystkie komendy z definicji, więc liczone dawałyby pełne pokrycie zawsze.
import json
import pathlib
import re
import sys

korzen = pathlib.Path(sys.argv[1]) if len(sys.argv) > 1 else pathlib.Path('.')
b = korzen / 'budowa'
kontrakt = json.loads((b / 'shared/contract.json').read_text(encoding='utf-8'))
wszystkie = {x['typ'] for x in kontrakt['komendy']}
stale = dict(re.findall(
    r"^\s*(\w+):\s*'([\w.]+)',",
    (b / 'shared/contract.ts').read_text(encoding='utf-8'),
    re.M,
))

wolane = set()
for p in (b / 'klient/src').rglob('*.ts'):
    if p.name in ('rejestr-komend.ts', 'contract.ts'):
        continue
    tresc_pliku = p.read_text(encoding='utf-8', errors='replace')
    for m in re.finditer(r'Command\.(\w+)', tresc_pliku):
        if m.group(1) in stale:
            wolane.add(stale[m.group(1)])

brak = sorted(wszystkie - wolane)
print(f'{len(wolane)}/{len(wszystkie)}')

if '--rodziny' in sys.argv:
    rodziny = {}
    for k in brak:
        rodziny.setdefault(k.split('.')[0], []).append(k)
    for r, lista in sorted(rodziny.items(), key=lambda x: -len(x[1])):
        w_rodzinie = len([k for k in wszystkie if k.split('.')[0] == r])
        print(f'{r} {len(lista)}/{w_rodzinie}')

if '--wykaz' in sys.argv:
    for k in brak:
        print(k)
