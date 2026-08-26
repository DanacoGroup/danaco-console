# Norma okna — zasada zera

> Końcowy efekt, który widzi użytkownik, nie może być efektem pracy w ostatecznym pliku.
> Okno **składa się** z biblioteki. Czegoś brakuje — najpierw trafia do biblioteki,
> dopiero potem jest załączane. To jest ta sama robota, a różnica jest taka, że następne
> okno dostaje ten składnik za darmo, zamiast wymyślać go po raz trzeci.
>
> — decyzja właściciela, sierpień 2026

Powód normy jest jeden i policzalny: aplikacja ma setki okien. Jeżeli każde okno wolno
dopisać choćby jedną barwę, to po stu oknach jest sto wyglądów. Spójności nie da się
utrzymać dyscypliną — da się ją utrzymać tylko tym, że nie ma gdzie jej złamać.

## Zasada

**W pliku okna i w arkuszu okna — zero wyglądu.**

| | należy do | wolno w oknie |
|---|---|---|
| barwa, tło, wypełnienie | biblioteka | nie |
| pismo: krój, stopień, grubość, interlinia, świat liter | biblioteka | nie |
| obramowanie, zaokrąglenie, cień, kontur, przezroczystość | biblioteka | nie |
| stany: najechanie, fokus, naciśnięcie, wybrany, nieczynny, błąd | biblioteka | nie |
| przejście, animacja | biblioteka | nie |
| siatka, kolumny, odstęp, kolejność, wyrównanie | — | **tak** |
| położenie, rozmiar, przewijanie, widoczność etapu | — | **tak** |

Składnik, który ma choć jedną deklarację z górnej części tabeli, nie jest rozmieszczeniem.
Jest komponentem i jego miejsce jest w `zasoby/css/komponenty.css`.

## Postępowanie, gdy czegoś brakuje

1. **Szukaj w bibliotece.** 261 składników. Ta sama rola bywa pod inną nazwą — sprawdzaj
   po roli i po deklaracjach, nie po nazwie.
2. **Nie znalazłeś — dopisz do biblioteki** jako `.dn-*`, z kompletem stanów
   (spoczynek · najechanie · naciśnięcie · fokus · wybrany · ładowanie · błąd)
   i wyłącznie na żetonach `--dn-*`.
3. **Dopiero teraz użyj go w oknie.**

Nigdy odwrotnie. Składnik dopisany „na razie w oknie, potem się przeniesie" nie zostaje
przeniesiony nigdy — zostaje skopiowany do następnego okna.

## Czego nie wolno nazywać brakiem

- **Wariantu.** Różnica między oknami wyraża się modyfikatorem `--*` przy istniejącym
  składniku, nie nowym składnikiem.
- **Rozmieszczenia.** Ten sam kafel w trzech kolumnach zamiast czterech to siatka okna,
  nie nowy kafel.
- **Roli, którą składa się z gotowych.** Nagłówek plus opis plus przycisk to trzy składniki
  biblioteki i jedna reguła siatki — nie czwarty składnik.

## Sprawdzenie

```
~/robocze/pomiar/straznik-skladnikow.sh
```

Strażnik odrzuca okno, które definiuje składnik o roli już obecnej w bibliotece,
oraz arkusz okna zawierający deklaracje wyglądu. Okno, które nie przechodzi strażnika,
nie jest gotowe — niezależnie od tego, jak wygląda na ekranie.

## Zakres

Norma obowiązuje wszystkie pliki w `05-okna/` oraz wszystkie arkusze w `zasoby/okna/`
i `zasoby/wejscie.css`. Wyjątkiem jest wyłącznie warstwa `.pt-*` — rusztowanie prototypu,
które nie wchodzi do produktu.

---

## Trzy strefy okna przedaplikacyjnego

> Pasek tytułu jest paskiem najciemniejszym. Lewa strona nie może być w takim
> samym kolorze co pasek tytułu, ani dolna belka nie może być w takim samym
> kolorze co pasek tytułu. To musi być rozdzielone i usztywnione, żeby było
> powtarzalne we wszystkich krokach i we wszystkich oknach tego typu.
>
> — decyzja właściciela, 26 sierpnia 2026

Okno przedaplikacyjne ma **trzy strefy i trzy różne barwy**, w obu motywach:

| strefa | żeton | co obejmuje |
|---|---|---|
| 1 · belka tytułowa — **najciemniejsza** | `--dn-okno-belka` | pasek tytułu okna |
| 2 · obrzeże | `--dn-okno-obrzeze` | kolumna lewa (kroki, tożsamość) **oraz pas działań na dole** |
| 3 · płótno | `--dn-okno-plotno` | panel treści kroku |

Wartości stoją w `zetony/zetony.css`, po jednym komplecie na motyw i na wybór
systemowy. **Okno nigdy nie ustala własnych** — bierze żeton.

Zmierzone rozdzielenie:

| motyw | belka ↔ obrzeże | obrzeże ↔ płótno |
|---|---:|---:|
| jasny | 87,5 L\* | 6,6 L\* |
| ciemny | 5,5 L\* | 8,8 L\* |

### Sprawdzenie

```
PLAYWRIGHT_BROWSERS_PATH=/opt/ms-playwright node ~/robocze/pomiar/strefy-okna.mjs
```

Kontrola przechodzi przez **każdy krok instalatora i każdy widok wejścia, w obu
motywach**, i odrzuca sytuację, w której belka ma barwę obrzeża, pas działań ma
barwę belki, obrzeże zlewa się z płótnem albo belka jest jaśniejsza od obrzeża.
