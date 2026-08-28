# Skrypt nakłada proponowane wpięcie na kopie dwóch plików pakietu models, dowodząc, że wpięcie się kompiluje, bez dotykania oryginalnych plików pakietu na dysku.
"""Nakłada proponowane wpięcie na kopie plików pakietu models.

Oryginałów nie dotyka: kopie leżą w katalogu roboczym, a kompilator widzi je
wyłącznie przez `go test -overlay`. Skrypt jest dowodem, że wpięcie naprawdę
się kompiluje.
"""
import pathlib
import sys

KAT = pathlib.Path(__file__).parent

WPIECIA = [
    (
        "adapter_api.go",
        # kotwica: import
        (
            '\t"net/http"\n\t"strings"\n\t"time"\n)',
            '\t"net/http"\n\t"strings"\n\t"time"\n\n'
            '\t"danacoconsole/server/internal/zewnetrzne"\n)',
        ),
    ),
    (
        "adapter_api.go",
        # kotwica: cialo bladOdpowiedzi
        (
            '\tsurowe, _ := io.ReadAll(io.LimitReader(odpowiedz.Body, 8<<10))\n'
            '\tkomunikat := strings.TrimSpace(string(surowe))\n'
            '\tif sciezka := d.Parametr("sciezka_bledu"); sciezka != "" {\n'
            '\t\tif wyjete, jest := WartoscZeSciezki(json.RawMessage(surowe), sciezka); jest {\n'
            '\t\t\tkomunikat = wyjete\n'
            '\t\t}\n'
            '\t}\n'
            '\treturn fmt.Errorf("models: kanał %s: odpowiedź %d: %s", d.Kod, odpowiedz.StatusCode, komunikat)',

            '\tsurowe, _ := io.ReadAll(io.LimitReader(odpowiedz.Body, 8<<10))\n'
            '\tsekret := zewnetrzne.SekretZOdpowiedzi(odpowiedz,\n'
            '\t\td.ParametrLub("naglowek_klucza", "Authorization"))\n'
            '\todmowa := zewnetrzne.NowaOdmowaKanaluZewnetrznego(d.Kod, d.PoswiadczenieOdwolanie,\n'
            '\t\todpowiedz.StatusCode, surowe, sekret)\n'
            '\tif sciezka := d.Parametr("sciezka_bledu"); sciezka != "" {\n'
            '\t\tif wyjete, jest := WartoscZeSciezki(json.RawMessage(surowe), sciezka); jest {\n'
            '\t\t\todmowa = odmowa.ZKomunikatem(wyjete, sekret)\n'
            '\t\t}\n'
            '\t}\n'
            '\treturn fmt.Errorf("models: %w", odmowa)',
        ),
    ),
    (
        "zadanie_api.go",
        (
            '\t"os"\n\t"strings"\n)',
            '\t"os"\n\t"strings"\n\n\t"danacoconsole/server/internal/zewnetrzne"\n)',
        ),
    ),
    (
        "zadanie_api.go",
        (
            '\tif byt, jestReferencja := strings.CutPrefix(odwolanie, przedrostekSejfu); jestReferencja {\n'
            '\t\tsejf := sejfPoswiadczen\n'
            '\t\tif sejf == nil {\n'
            '\t\t\treturn "", fmt.Errorf("models: kanał %s: odwołanie %s wskazuje sejf, a sejfu nie wpięto", d.Kod, odwolanie)\n'
            '\t\t}\n'
            '\t\tklucz, jest := sejf.Odczytaj(ctx, strings.TrimSpace(byt))\n'
            '\t\tif !jest || strings.TrimSpace(klucz) == "" {\n'
            '\t\t\treturn "", fmt.Errorf("models: kanał %s: brak poświadczenia w sejfie pod odwołaniem %s", d.Kod, odwolanie)\n'
            '\t\t}\n'
            '\t\treturn klucz, nil\n'
            '\t}\n'
            '\tklucz := strings.TrimSpace(os.Getenv(odwolanie))\n'
            '\tif klucz == "" {\n'
            '\t\treturn "", fmt.Errorf("models: kanał %s: brak wartości pod odwołaniem %s", d.Kod, odwolanie)\n'
            '\t}\n'
            '\treturn klucz, nil',

            '\topis := zewnetrzne.RozpoznajPoswiadczenie(odwolanie)\n'
            '\tif opis.ZSejfu() {\n'
            '\t\tsejf := sejfPoswiadczen\n'
            '\t\tif sejf == nil {\n'
            '\t\t\treturn "", zewnetrzne.NowyBrakPoswiadczenia(d.Kod, odwolanie, zewnetrzne.PowodSejfNiewpiety)\n'
            '\t\t}\n'
            '\t\tklucz, jest := sejf.Odczytaj(ctx, opis.Byt)\n'
            '\t\tif !jest {\n'
            '\t\t\treturn "", zewnetrzne.NowyBrakPoswiadczenia(d.Kod, odwolanie, zewnetrzne.PowodBrakWpisu)\n'
            '\t\t}\n'
            '\t\tif strings.TrimSpace(klucz) == "" {\n'
            '\t\t\treturn "", zewnetrzne.NowyBrakPoswiadczenia(d.Kod, odwolanie, zewnetrzne.PowodPustyWpis)\n'
            '\t\t}\n'
            '\t\treturn klucz, nil\n'
            '\t}\n'
            '\tklucz := strings.TrimSpace(os.Getenv(odwolanie))\n'
            '\tif klucz == "" {\n'
            '\t\treturn "", zewnetrzne.NowyBrakPoswiadczenia(d.Kod, odwolanie, zewnetrzne.PowodPustaZmienna)\n'
            '\t}\n'
            '\treturn klucz, nil',
        ),
    ),
]


def main() -> int:
    for plik, (kotwica, tresc) in WPIECIA:
        sciezka = KAT / plik
        zrodlo = sciezka.read_text(encoding="utf-8")
        if zrodlo.count(kotwica) != 1:
            print(f"KOTWICA NIETRAFIONA w {plik}: wystąpień {zrodlo.count(kotwica)}")
            return 1
        sciezka.write_text(zrodlo.replace(kotwica, tresc), encoding="utf-8")
        print(f"wpięto: {plik}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
