// Wycięcie szkieletu witryny grupy z jej ŻYWEJ strony — i złożenie z niego dwóch
// szablonów: dla paczki do wgrania oraz dla kanału wydań.
//
// PO CO TO POLECENIE, A NIE SKOPIOWANY RAZ KAWAŁEK HTML-A. Bo nagłówek i stopka
// witryny grupy są jej własnością i będą się zmieniać bez naszego udziału.
// Przepisany ręcznie nagłówek zestarzeje się w tygodniu, w którym ktoś doda tam
// pozycję menu — i wtedy strona o produkcie stanie się jedyną podstroną z innym
// menu niż pozostałe. Szkielet wycina się więc poleceniem z pobranej strony
// (`curl https://www.danaco-group.pl/obszary.html`), a nie kopiuje wzrokiem.
//
//   node portfolio/wytnij-szkielet.mjs <sciezka-do-obszary.html>
//
// Powstają dwa szablony, bo to dwa różne miejsca w sieci:
//
//   szkielet-grupy.html   — dla paczki wgrywanej NA danaco-group.pl. Odnośniki
//                           zostają względne, bo strona stanie obok pozostałych.
//   szkielet-kanalu.html  — dla kanału wydań (pobierz.danaco-group.pl). To OSOBNY
//                           host: `o-nas.html` nie istnieje tam i dałoby błąd 404,
//                           więc odnośniki do stron grupy są rozwinięte do adresu
//                           bezwzględnego — podobnie arkusz, godło i favikon,
//                           które są LINKOWANE z witryny grupy, nie kopiowane:
//                           kopia rozjechałaby się przy pierwszej ich poprawce
//                           i kanał zostałby jedyną stroną w starym wyglądzie.
//
// TRZY POZYCJE MENU DOPISYWANE SĄ TUTAJ, nie ręką w pliku wynikowym — bo mają
// wejść w oba szablony jednakowo, a to trzy różne miejsca w tym samym markupie.
import { readFile, writeFile } from 'node:fs/promises';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

const KATALOG = dirname(fileURLToPath(import.meta.url));
const zrodloSciezka = process.argv[2];
if (!zrodloSciezka) {
  console.error(
    'brak ścieżki do pobranej strony witryny grupy.\n' +
      '  użycie: node portfolio/wytnij-szkielet.mjs /sciezka/obszary.html\n' +
      '  stronę pobiera się: curl -s -o obszary.html https://www.danaco-group.pl/obszary.html',
  );
  process.exit(1);
}

// Znacznik kolejności bajtów i końce wierszy CRLF — żywa witryna niesie oba.
// Bez normalizacji kotwice menu (dopasowywane całym wierszem) nigdy się nie
// trafiają, a polecenie odmawia z komunikatem o „zmienionym układzie menu",
// który jest nieprawdą: układ jest ten sam, różnią się niewidzialne bajty.
const zrodlo = (await readFile(zrodloSciezka, 'utf8')).replace(/^﻿/, '').replaceAll('\r\n', '\n');

function wytnij(od, doKonca) {
  const a = zrodlo.indexOf(od);
  if (a < 0) throw new Error(`nie znaleziono w stronie źródłowej: ${od}`);
  if (doKonca === null) return zrodlo.slice(a);
  const b = zrodlo.indexOf(doKonca, a);
  if (b < 0) throw new Error(`nie znaleziono w stronie źródłowej: ${doKonca}`);
  return zrodlo.slice(a, b + doKonca.length);
}

let naglowek = wytnij('<a class="skip-link"', '</header>');
let stopka = wytnij('<footer class="site-footer">', null);

// --- Trzy pozycje menu ------------------------------------------------------
// Kotwice są WIERSZAMI ISTNIEJĄCEGO menu, nie odległościami w znakach: gdy tamta
// witryna przestawi kolejność, polecenie padnie z nazwą brakującego wiersza,
// zamiast wstawić pozycję w przypadkowe miejsce.
function wstawPrzed(gdzie, kotwica, wiersz, nazwa) {
  if (!gdzie.includes(kotwica)) {
    throw new Error(
      `nie znaleziono miejsca wstawienia pozycji „${nazwa}". Szukany wiersz menu:\n  ${kotwica.trim()}\n` +
        '  Menu witryny grupy zmieniło układ — trzeba wskazać nową kotwicę, a nie wstawiać na oślep.',
    );
  }
  return gdzie.replace(kotwica, wiersz + kotwica);
}

const POZYCJE = [
  {
    nazwa: 'Portfolio → Danaco Console',
    kotwica: '              <a class="nav__dd-link" href="obszary.html#hotelarstwo" role="menuitem">Nasze marki hotelowe</a>\n',
    wiersz: '              <a class="nav__dd-link" href="danaco-console.html" role="menuitem">Danaco Console</a>\n',
  },
  {
    nazwa: 'Danaco Share → Danaco Console',
    kotwica: '              <a class="nav__dd-link" href="dokumenty.html" role="menuitem">Repozytorium dokumentów</a>\n',
    wiersz: '              <a class="nav__dd-link" href="danaco-console.html#platforma" role="menuitem">Danaco Console</a>\n',
  },
  {
    // Pozycja paska górnego — BEZ rozwinięcia. Cztery sąsiednie pozycje mają
    // podmenu, ta nie ma czego w nim trzymać: prowadzi wprost do plików.
    // Bez `target="_blank"`: to nadal witryna grupy, tylko na subdomenie, a
    // otwieranie nowej karty przy przejściu w obrębie własnych stron jest natrętne.
    nazwa: 'pasek górny → Pobierz',
    kotwica: '      <a class="nav__link" href="kontakt.html">Kontakt</a>',
    wiersz: '      <a class="nav__link" href="https://pobierz.danaco-group.pl/" rel="noopener">Pobierz</a>\n',
  },
];

for (const p of POZYCJE) naglowek = wstawPrzed(naglowek, p.kotwica, p.wiersz, p.nazwa);

// --- Szablon dla kanału wydań: odnośniki do stron grupy rozwinięte ----------
// Rozwijane są odnośniki do stron grupy ORAZ do jej katalogu `assets/`. Arkusz
// witryny grupy jest LINKOWANY, nie kopiowany — jego kopia na kanale rozjechałaby
// się przy pierwszej ich poprawce, a wtedy jedna podstrona wyglądałaby inaczej
// niż cała witryna i nikt by nie wiedział dlaczego. Wersja arkusza (v=20260602f)
// jest tym, wobec czego strona była składana; przy ich zmianie wersji trzeba
// stronę przejrzeć, a nie zakładać, że nadal pasuje.
// Nie rozwijamy odnośników już bezwzględnych, zakotwiczeń ani `mailto:`.
const GRUPA = 'https://www.danaco-group.pl/';
function rozwinOdnosniki(kod) {
  return kod.replace(/(href|src)="(?!https?:|mailto:|#)([^"]+)"/g, (calosc, atrybut, adres) => {
    if (adres === 'danaco-console.html' || adres.startsWith('danaco-console.html#')) {
      return `${atrybut}="${GRUPA}${adres}"`;
    }
    if (adres.endsWith('.html') || adres.startsWith('dokumenty/')) return `${atrybut}="${GRUPA}${adres}"`;
    if (adres.startsWith('assets/')) return `${atrybut}="${GRUPA}${adres}"`;
    return calosc;
  });
}

function szablon({ kanoniczny, naglowekKodu, stopkaKodu, podnawigacja, zasoby }) {
  return `<!DOCTYPE html>
<html lang="pl">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{TYTUL}}</title>
<meta name="description" content="{{OPIS}}">
<meta name="theme-color" content="#0B1322">
<link rel="icon" href="${zasoby}assets/img/favicon.png" type="image/png">
<link rel="canonical" href="${kanoniczny}">
<link rel="preconnect" href="https://fonts.googleapis.com">
<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
<link href="https://fonts.googleapis.com/css2?family=DM+Sans:wght@300;400;500;700;800&family=Plus+Jakarta+Sans:wght@300;400;500;600;700&family=Manrope:wght@300;400;500;600;700;800&family=Fraunces:ital,opsz,wght@0,9..144,400;0,9..144,500;1,9..144,400;1,9..144,500&display=swap" rel="stylesheet">
<link href="https://api.fontshare.com/v2/css?f[]=general-sans@200,300,400,500,600,700&display=swap" rel="stylesheet">
<meta name="robots" content="index,follow,max-image-preview:large">
<meta property="og:type" content="website">
<meta property="og:site_name" content="Danaco Group">
<meta property="og:title" content="{{TYTUL}}">
<meta property="og:description" content="{{OPIS}}">
<meta property="og:locale" content="pl_PL">
<link rel="stylesheet" href="${zasoby}assets/css/styles.css?v=20260602f">
{{STYL}}
</head>
<body>
${naglowekKodu}
${podnawigacja}
<main id="main">
{{TRESC}}
</main>

${stopkaKodu}`;
}

// Podnawigacja stoi TYLKO na kanale wydań. Kanał ma dziesięć własnych stron
// (Pobierz, Prototypy, Moduły…), których menu witryny grupy nie zna i nie ma
// powodu znać. Bez tego paska te strony byłyby osiągalne wyłącznie odnośnikami
// w treści — czyli po omacku.
const PODNAWIGACJA = `<nav class="dc-podnaw" aria-label="Strony o produkcie">
  <div class="container">
    <span class="dc-podnaw__ttl">Danaco Console</span>
{{TRASY}}
  </div>
</nav>`;

await writeFile(
  join(KATALOG, 'szkielet-grupy.html'),
  szablon({
    kanoniczny: `${GRUPA}danaco-console.html`,
    naglowekKodu: naglowek,
    stopkaKodu: stopka,
    podnawigacja: '',
    zasoby: '',
  }),
  'utf8',
);

await writeFile(
  join(KATALOG, 'szkielet-kanalu.html'),
  szablon({
    kanoniczny: 'https://pobierz.danaco-group.pl/{{ADRES}}',
    naglowekKodu: rozwinOdnosniki(naglowek),
    stopkaKodu: rozwinOdnosniki(stopka),
    podnawigacja: PODNAWIGACJA,
    zasoby: GRUPA,
  }),
  'utf8',
);

console.log('wycięte z:', zrodloSciezka);
console.log('szablony: portfolio/szkielet-grupy.html, portfolio/szkielet-kanalu.html');
console.log('dopisane pozycje menu:', POZYCJE.map((p) => p.nazwa).join(' · '));
