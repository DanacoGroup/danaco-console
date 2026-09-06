# Wkład w repozytorium

Danaco Console jest oprogramowaniem własnościowym Danaco Holding Group.
Kod jest widoczny publicznie, ale repozytorium nie przyjmuje wkładu z zewnątrz:
propozycje zmian i rewizje od osób spoza zespołu nie są scalane.

Warunki korzystania określa [docs/LICENSE.md](docs/LICENSE.md). Sama widoczność
kodu nie udziela licencji na użycie, zwielokrotnianie ani tworzenie opracowań.

Wadę bezpieczeństwa zgłasza się drogą opisaną w [SECURITY.md](SECURITY.md).
Pozostałe sprawy — zwykłą drogą wsparcia: support@danaco-group.pl.

## Dla zespołu

Rewizja wchodzi przez drabinę weryfikacji:

```bash
git config core.hooksPath .githooks
bash narzedzia/drabina.sh szybka
```

Na GitHubie te same szczeble powtarza przepływ `Kontrola`
([.github/workflows/kontrola.yml](.github/workflows/kontrola.yml)).
Wytworów generatora kontraktu nie zmienia się ręcznie — powstają
z `budowa/shared/contract.json`.

---

© Danaco Holding Group. Wszelkie prawa zastrzeżone.
