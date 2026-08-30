# Miara wpięcia interfejsu: ile komend kontraktu woła klient.
import json, pathlib, re, sys

korzen = pathlib.Path(sys.argv[1])
kontrakt = json.loads((korzen / 'budowa/shared/contract.json').read_text(encoding='utf-8'))
wszystkie = {k['typ'] for k in kontrakt['komendy']}

stale = dict(re.findall(
    r"^\s*(\w+):\s*'([\w.]+)',",
    (korzen / 'budowa/shared/contract.ts').read_text(encoding='utf-8'),
    re.M,
))

wolane = set()
for plik in (korzen / 'budowa/klient/src').rglob('*.ts'):
    for m in re.finditer(r'Command\.(\w+)', plik.read_text(encoding='utf-8')):
        if m.group(1) in stale:
            wolane.add(stale[m.group(1)])

print(f'{len(wolane)}/{len(wszystkie)}')
