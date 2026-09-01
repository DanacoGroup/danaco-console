// Pomiar plików wydania — wpisuje do `wydania.json` rozmiar i sumę SHA-256
// odczytane Z PLIKU NA DYSKU.
//
// PO CO OSOBNE POLECENIE. Bo suma przepisana ręcznie jest sumą przepisaną
// z czegoś: z cudzego raportu, z pliku .sha256 obok, z poprzedniego przebiegu.
// Każde z tych źródeł bywa starsze niż plik. 18.08.2026 kosztowało to jeden
// nieprawdziwy wykaz: trzy pakiety przebudowano W TRAKCIE składania witryny
// i strona zdążyła podać sumy plików, których już nie było — a suma jest tu
// JEDYNYM sprawdzianem, jaki ma pobierający. Wykaz mierzy się więc poleceniem,
// nie ręką, i mierzy się PO zakończeniu budowy, nie w jej trakcie.
//
//   node zmierz.mjs            — mierzy i zapisuje
//   node zmierz.mjs --sprawdz  — tylko mówi, co się rozjechało; nic nie zapisuje
//
// Po pomiarze składa się witrynę ponownie: `node zloz.mjs`.
import { createReadStream } from 'node:fs';
import { readFile, writeFile, stat } from 'node:fs/promises';
import { createHash } from 'node:crypto';
import { dirname, join, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const KATALOG = dirname(fileURLToPath(import.meta.url));
const WYKAZ = join(KATALOG, 'wydania.json');
const TYLKO_SPRAWDZ = process.argv.includes('--sprawdz');

const wykaz = JSON.parse(await readFile(WYKAZ, 'utf8'));

// Skąd brać pliki. Ścieżka stoi w wykazie, a nie tutaj, bo katalog wydania zmienia
// się z każdą datą — a to polecenie ma zostać takie samo.
const katalogPlikow = resolve(KATALOG, wykaz.katalog_plikow ?? '../wydania');

/** SHA-256 liczone strumieniem: pakiet natywny waży setki megabajtów i wciąganie
 *  go w całości do pamięci byłoby zapłatą za nic. */
async function suma(sciezka) {
  const skrot = createHash('sha256');
  for await (const kawalek of createReadStream(sciezka)) skrot.update(kawalek);
  return skrot.digest('hex');
}

/** Zapis rozmiaru dla człowieka — mebibajty, przecinek po polsku. */
function rozmiarLudzki(bajty) {
  return `${(bajty / 1048576).toFixed(1).replace('.', ',')} MB`;
}

/** Wszystkie pozycje niosące `nazwaPliku`, wraz z nazwą miejsca, w którym leżą —
 *  nazwa miejsca wchodzi do meldunku, żeby było wiadomo, co dokładnie się zmieniło. */
function pozycje(w) {
  const zebrane = [];
  for (const [klucz, wartosc] of Object.entries(w)) {
    // `pola` to słownik OPISUJĄCY pola wykazu, nie pozycja wydania — a niesie
    // klucz `nazwaPliku` z opisem w wartości i bez tego wyjątku polecenie
    // meldowałoby przy każdym przebiegu brak pliku o nazwie „nazwa pokazywana
    // na stronie". Meldunek, który zawsze kłamie, uczy nie czytać meldunków.
    // `w_przygotowaniu` z tego samego powodu: pozycja zapowiedziana nie ma
    // jeszcze pliku i to jest jej stan, a nie rozjazd do zgłoszenia.
    if (klucz === 'pola' || klucz === 'w_przygotowaniu') continue;
    if (Array.isArray(wartosc)) {
      wartosc.forEach((p, i) => {
        if (p && typeof p === 'object' && p.nazwaPliku) zebrane.push({ p, gdzie: `${klucz}[${i}]` });
      });
    } else if (wartosc && typeof wartosc === 'object' && wartosc.nazwaPliku) {
      zebrane.push({ p: wartosc, gdzie: klucz });
    }
  }
  return zebrane;
}

let zmian = 0;
let brakow = 0;

for (const { p, gdzie } of pozycje(wykaz)) {
  const sciezka = join(katalogPlikow, p.nazwaPliku);
  let dane;
  try {
    dane = await stat(sciezka);
  } catch {
    // Brak pliku jest BRAKIEM, nie stanem przejściowym: pozycja zapowiedziana
    // stoi w `w_przygotowaniu` i tu jej nie ma. Pozycja publikowana, której pliku
    // nie ma na dysku, niesie sumę wziętą skądinąd — a suma nie z pliku jest
    // dokładnie tym, przed czym stoi to polecenie.
    console.log(`— ${gdzie}: pliku nie ma na dysku (${p.nazwaPliku})`);
    brakow += 1;
    continue;
  }
  const policzona = await suma(sciezka);
  const rozjazd = [];
  if (p.rozmiarBajty !== dane.size) rozjazd.push(`rozmiar ${p.rozmiarBajty ?? '—'} → ${dane.size}`);
  if (p.suma !== policzona) rozjazd.push(`suma ${(p.suma ?? '—').slice(0, 12)}… → ${policzona.slice(0, 12)}…`);

  if (rozjazd.length === 0) {
    console.log(`= ${gdzie}: zgadza się (${p.nazwaPliku})`);
    continue;
  }
  zmian += 1;
  console.log(`${TYLKO_SPRAWDZ ? '!' : '~'} ${gdzie}: ${rozjazd.join(', ')}`);
  if (TYLKO_SPRAWDZ) continue;
  p.rozmiarBajty = dane.size;
  p.rozmiar = rozmiarLudzki(dane.size);
  p.suma = policzona;
}

if (!TYLKO_SPRAWDZ && zmian > 0) {
  await writeFile(WYKAZ, `${JSON.stringify(wykaz, null, 2)}\n`, 'utf8');
}

console.log(
  TYLKO_SPRAWDZ
    ? `sprawdzone: rozjazdów ${zmian}, plików brakuje ${brakow} — nic nie zapisano`
    : `zmierzone: poprawionych pozycji ${zmian}, plików brakuje ${brakow}` +
        (zmian > 0 ? ' — złóż witrynę ponownie: node zloz.mjs' : ''),
);

// Rozjazd i brak są przy `--sprawdz` odmową, nie uwagą: to polecenie ma dać się
// wpiąć przed wgraniem i zatrzymać wgranie wykazu, który mówi o innych plikach
// niż te leżące na dysku — albo o plikach, których na dysku nie ma wcale i nie
// ma z czego policzyć ich sumy. Wykaz publikuje się po przebiegu, który mówi
// „rozjazdów 0, plików brakuje 0".
if (TYLKO_SPRAWDZ && (zmian > 0 || brakow > 0)) process.exit(1);
