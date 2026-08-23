// Strona „Prototypy" — żywe okna produktu, nie obrazki.
//
// DLACZEGO RAMKA, A NIE ZRZUT. Makiety z dostawy (`design/05-okna/`) są
// interaktywne: każda niesie własne skrypty i wspólny `zasoby/prototyp.js`.
// Zrzut ekranu zabiera z nich dokładnie to, co jest w nich najlepsze — przejścia
// zakładek, rozwijane panele, przełączniki. Dlatego okna wchodzą tu w ramkach
// (`<iframe>`), a nie przez wklejenie ich HTML-a w treść strony: mają własną
// warstwę stylów i własny skrypt, które po wklejeniu zderzyłyby się ze stroną
// (te same nazwy klas, ten sam `document`).
//
// UKŁAD KATALOGÓW MUSI ZOSTAĆ ZACHOWANY. Makiety wołają zasoby ścieżkami
// względnymi w rodzaju `../../zasoby/css/fundament.css`. Kopiowanie samych plików
// HTML dałoby okna bez stylów. Dlatego `zloz.mjs` przenosi do wyjścia CAŁE
// `05-okna/` wraz z CAŁYM `zasoby/`, zachowując wzajemne położenie — pod
// `dist/prototypy/`, gdzie `../../` z `05-okna/moduly/okno.html` trafia
// w `prototypy/zasoby/`.
//
// BEZ JAVASCRIPTU STRONA DZIAŁA. Każde okno jest zwykłym odnośnikiem do pliku
// makiety i otworzy się kliknięciem także przy wyłączonym skrypcie. Skrypt
// jedynie PODMIENIA zawartość ramki zamiast opuszczać stronę. Nic, co Operator
// ma zobaczyć, nie powstaje dopiero ze skryptu.
import { readdir, readFile } from 'node:fs/promises';
import { join } from 'node:path';

import { tekst } from './szkielet.mjs';

// Nazwy rodzin nie są wymyślone — to nazwy katalogów dostawy. Opis obok mówi,
// co w danej rodzinie stoi, i nie dodaje ani jednego okna ponad te, które są.
// Pole `zrodlo` jest ścieżką w dostawie i jednocześnie ścieżką w wyjściu (pod
// `prototypy/`) — bo układ katalogów przenosi się bez zmian, inaczej odwołania
// `../../zasoby/…` z wnętrza makiet przestałyby trafiać.
const RODZINY = [
  {
    zrodlo: '_samodzielne',
    nazwa: 'Pełne stanowiska',
    opis: 'Prototypy samodzielne — cały interfejs w jednym pliku, bez odwołań na zewnątrz.',
  },
  { zrodlo: '05-okna/srodowiska', nazwa: 'Środowiska', opis: 'Cztery przestrzenie pracy i powłoka orkiestracji.' },
  { zrodlo: '05-okna/moduly', nazwa: 'Moduły', opis: 'Okna robocze poszczególnych modułów.' },
  { zrodlo: '05-okna/przeplyw', nazwa: 'Przepływ', opis: 'Droga Operatora od uruchomienia do strefy roboczej.' },
  {
    zrodlo: '05-okna/platformowe',
    nazwa: 'Platformowe',
    opis: 'Okna wspólne dla całej platformy — ustawienia, konfiguracja, postać mobilna, ekran zawsze widoczny.',
  },
  { zrodlo: '05-okna', nazwa: 'Wzorzec', opis: 'Wzorzec stanowiska — układ, z którego wyprowadzone są okna modułów.' },
];

/** Tytuł okna czytany Z MAKIETY, nie nadany przez wykonawcę. Kanon zabrania
 *  wymyślania nazw okien — więc nazwa bierze się ze znacznika `<title>` samej
 *  makiety, po odjęciu powtarzającego się przedrostka z nazwą produktu. */
async function tytulMakiety(sciezka) {
  const zrodlo = await readFile(sciezka, 'utf8');
  const dopasowanie = zrodlo.match(/<title>([\s\S]*?)<\/title>/i);
  if (!dopasowanie) return null;
  return dopasowanie[1]
    .replace(/\s+/g, ' ')
    .trim()
    .replace(/^Danaco Console\s*[—–-]\s*/, '');
}

/** Wykaz okien jednej rodziny, ułożony nazwą pliku — ta sama kolejność, co
 *  w katalogu dostawy, więc numerowane okna przepływu stoją po kolei. */
async function oknaRodziny(korzen, zrodlo) {
  let pliki;
  try {
    // `withFileTypes` odsiewa katalogi: rodzina „Wzorzec" wskazuje `05-okna`,
    // w którym poza wzorcem leżą katalogi pozostałych rodzin.
    pliki = (await readdir(join(korzen, zrodlo), { withFileTypes: true }))
      .filter((p) => p.isFile() && p.name.endsWith('.html'))
      .map((p) => p.name)
      .sort();
  } catch {
    return [];
  }
  const okna = [];
  for (const plik of pliki) {
    const tytul = await tytulMakiety(join(korzen, zrodlo, plik));
    okna.push({
      plik,
      adres: `prototypy/${zrodlo}/${plik}`,
      tytul: tytul ?? plik.replace(/\.html$/, ''),
      // POSTAĆ MOBILNA I EKRAN ZAWSZE WIDOCZNY — CO TU NAPRAWDĘ JEST.
      //
      // Pierwsze podejście wstawiało te dwie makiety w ramkę o proporcjach
      // telefonu (520 px szerokości), bo tego dotyczą. Wyszło źle i pomiar
      // powiedział, dlaczego: `mobile.html` NIE JEST dokumentem o szerokości
      // telefonu. To strona opisowa o szerokości do 1100 px, która POKAZUJE
      // W ŚRODKU ekran telefonu (348 px) wraz z opisem funkcji globalnej —
      // tak samo `always-on-display.html`. Wciśnięta w 520 px rozjeżdżała się,
      // bo ma własny próg 1100 px i poniżej niego przestawia układ.
      //
      // Stąd rozstrzygnięcie: obie idą w szerokości WŁASNEJ (1100 px), w osobnej
      // sekcji „Postać mobilna" i jedna pod drugą, bo są wyższe niż szersze.
      // Ekran telefonu widać w nich taki, jaki zaprojektowano — tylko nie jako
      // ramkę strony, a jako zawartość makiety. Zmyślanie ramki telefonu wokół
      // dokumentu, który jej nie ma, byłoby pokazaniem produktu w układzie,
      // którego dostawa nie przewiduje.
      wlasna: plik === 'mobile.html' || plik === 'always-on-display.html',
    });
  }
  return okna;
}

/** Zbiera cały katalog makiet. Zwraca też liczbę okien — wchodzi do treści
 *  strony, żeby nikt nie musiał jej liczyć ręcznie i się pomylić. */
export async function katalogPrototypow(korzen) {
  const rodziny = [];
  for (const rodzina of RODZINY) {
    const okna = await oknaRodziny(korzen, rodzina.zrodlo);
    if (okna.length > 0) rodziny.push({ ...rodzina, okna, klucz: rodzina.zrodlo.replace(/[^a-z0-9]+/gi, '-') });
  }
  return rodziny;
}

const SKRYPT = `
(function () {
  var ramka = document.getElementById('ramka-prototypu');
  var podpis = document.getElementById('podpis-prototypu');
  var pelne = document.getElementById('pelne-okno');
  if (!ramka || !podpis || !pelne) return;

  function wybierz(odnosnik) {
    var adres = odnosnik.getAttribute('href');
    var nazwa = odnosnik.getAttribute('data-tytul') || odnosnik.textContent.trim();
    ramka.setAttribute('src', adres);
    ramka.setAttribute('title', 'Prototyp okna: ' + nazwa);
    podpis.textContent = nazwa;
    pelne.setAttribute('href', adres);
    var wszystkie = document.querySelectorAll('.dc-okno');
    for (var i = 0; i < wszystkie.length; i += 1) {
      wszystkie[i].removeAttribute('aria-current');
    }
    odnosnik.setAttribute('aria-current', 'true');
    // Ramka jest niżej niż wykaz na wąskim ekranie — bez tego kliknięcie
    // wyglądałoby, jakby nic się nie stało.
    if (window.matchMedia('(max-width: 960px)').matches) {
      ramka.scrollIntoView({ block: 'start', behavior: 'smooth' });
    }
  }

  document.addEventListener('click', function (zdarzenie) {
    var odnosnik = zdarzenie.target.closest ? zdarzenie.target.closest('.dc-okno') : null;
    if (!odnosnik) return;
    // Nowa karta, zapis pliku, środkowy przycisk — zostawiamy przeglądarce.
    if (zdarzenie.metaKey || zdarzenie.ctrlKey || zdarzenie.shiftKey || zdarzenie.button !== 0) return;
    zdarzenie.preventDefault();
    wybierz(odnosnik);
  });

  var rodziny = document.querySelectorAll('.dc-rodzina');
  function pokazRodzine(klucz) {
    for (var i = 0; i < rodziny.length; i += 1) {
      var przycisk = rodziny[i];
      var wlasna = przycisk.getAttribute('data-rodzina') === klucz;
      przycisk.setAttribute('aria-selected', wlasna ? 'true' : 'false');
    }
    var wykazy = document.querySelectorAll('.dc-okna');
    for (var j = 0; j < wykazy.length; j += 1) {
      wykazy[j].hidden = wykazy[j].getAttribute('data-rodzina') !== klucz;
    }
  }
  for (var k = 0; k < rodziny.length; k += 1) {
    rodziny[k].addEventListener('click', function () {
      pokazRodzine(this.getAttribute('data-rodzina'));
    });
  }
  if (rodziny.length > 0) pokazRodzine(rodziny[0].getAttribute('data-rodzina'));

  // Dopasowanie skali: okno wyrenderowane w szerokości projektowej zostaje
  // pomniejszone dokładnie tyle, ile trzeba, żeby zmieściło się w scenie.
  // Wysokość sceny idzie za skalą, więc pod ramką nie zostaje pusty pas.
  function dopasuj() {
    var sceny = document.querySelectorAll('.dc-scena');
    for (var i = 0; i < sceny.length; i += 1) {
      var scena = sceny[i];
      var ramka = scena.querySelector('.dc-ramka');
      if (!ramka) continue;
      var szer = parseInt(scena.getAttribute('data-szer'), 10) || 1440;
      var wys = parseInt(scena.getAttribute('data-wys'), 10) || 900;
      var skala = scena.clientWidth / szer;
      if (!isFinite(skala) || skala <= 0) continue;
      ramka.style.transform = 'scale(' + skala + ')';
      scena.style.height = Math.round(wys * skala) + 'px';
    }
  }
  dopasuj();
  window.addEventListener('resize', dopasuj);
  // Zmiana okna w ramce nie zmienia jej rozmiaru, ale zmiana rodziny może
  // przestawić układ strony — przeliczamy po każdym kliknięciu w wykaz.
  document.addEventListener('click', function () { window.setTimeout(dopasuj, 0); });
})();
`;

const STYL = `
/* Reguły układu strony „Prototypy" — w KANONIE GRUPY, bo to strona na
   danaco-group.pl. Ani jednego żetonu --dn- poza wnętrzem ramki: barwy,
   odstępy, promienie i kroje biorą się z arkusza witryny grupy. Progi widoku
   idą po max-width, bo taką konwencję ma tamten arkusz (1100, 860, 760, 600 px)
   — kanon Console liczy je po min-width i mieszanie obu w jednym pliku dałoby
   dwa różne języki układu na jednej stronie. */
.dc-prototyp { display: grid; gap: var(--space-lg); grid-template-columns: 17rem 1fr; align-items: start; margin: var(--space-lg) 0; }
.dc-prototyp__wybor { position: sticky; top: 96px; }
.dc-rodziny { display: flex; flex-wrap: wrap; gap: 8px; margin-bottom: var(--space-sm); }
.dc-rodzina {
  font-family: var(--ff-body); font-size: .8rem; letter-spacing: .04em;
  min-height: 34px; padding: 0 14px; cursor: pointer;
  background: transparent; color: var(--text-muted);
  border: 1px solid var(--border-default); border-radius: var(--radius-pill);
}
.dc-rodzina:hover { color: var(--text-strong); border-color: var(--border-strong); }
/* Wybrana rodzina nie jest oznaczona samym kolorem — dochodzi grubość napisu. */
.dc-rodzina[aria-selected="true"] { background: var(--bg-inverse); color: var(--bg-surface); border-color: var(--bg-inverse); font-weight: 600; }
.dc-okna { list-style: none; margin: 0; padding: 0; border: 1px solid var(--border-default); border-radius: var(--radius-lg); overflow: hidden; background: var(--bg-surface); }
.dc-okna li + li { border-top: 1px solid var(--border-default); }
.dc-okno { display: block; padding: 10px 14px; color: var(--text-body); text-decoration: none; font-size: .86rem; line-height: 1.35; }
.dc-okno:hover { background: var(--bg-subtle); }
.dc-okno[aria-current="true"] { background: var(--bg-subtle); font-weight: 600; box-shadow: inset 3px 0 0 var(--brand-yellow); }
.dc-okno small { display: block; color: var(--text-muted); font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace; font-size: .72rem; }
.dc-scena__pasek { display: flex; align-items: center; gap: var(--space-sm); flex-wrap: wrap; margin-bottom: 8px; }
.dc-scena__podpis { font-family: var(--ff-display); font-size: .74rem; letter-spacing: .08em; text-transform: uppercase; color: var(--text-muted); }
.dc-scena__pelne { margin-left: auto; font-size: .86rem; }

/* OKNO RENDERUJE SIĘ W SWOJEJ SZEROKOŚCI I DOPIERO POTEM JEST POMNIEJSZANE.
   Makiety biurkowe projektowane są pod szeroki ekran. Wstawione do ramki szerokiej
   na 800 px pokazywały prawą krawędź uciętą w pół przycisku — czyli produkt
   w układzie, dla którego nie był robiony. Dlatego ramka ma STAŁĄ szerokość
   projektową i jest przeskalowana do miejsca, jakie ma na stronie. Skala liczona
   jest skryptem, bo CSS nie umie wyrazić „szerokość rodzica podzielona przez
   1440"; bez skryptu ramka zostaje w rozmiarze projektowym i przewija się
   w poziomie — mniej wygodnie, ale nadal całe okno.

   WNĘTRZE RAMKI JEST W KANONIE CONSOLE i tak ma być: pokazuje produkt, nie stronę.
   Ramka jest granicą obu kanonów — poza nią arkusz grupy, wewnątrz fundament.css
   i prototyp.css dostawy. Dwa drzewa stylów nie mają jak się zmieszać. */
.dc-scena { position: relative; overflow: hidden; border: 1px solid var(--border-default); border-radius: var(--radius-lg); background: var(--bg-surface); height: 620px; }
.dc-scena--wysoka { height: 860px; }
.dc-ramka { position: absolute; top: 0; left: 0; width: var(--szer, 1440px); height: var(--wys, 900px); border: 0; transform-origin: 0 0; display: block; }
@media (max-width: 1100px) {
  .dc-prototyp { grid-template-columns: 1fr; }
  .dc-prototyp__wybor { position: static; }
  /* Na telefonie okno biurkowe pomniejszone do 360 px byłoby nieczytelne, więc
     scena dostaje wysokość liczoną ekranem czytającego — a najważniejszy staje
     się wtedy odnośnik „otwórz w pełnym oknie" obok. */
  .dc-scena, .dc-scena--wysoka { height: 60vh; min-height: 300px; }
}
`;

/** Buduje treść strony „Prototypy" z katalogu makiet. */
export function stronaPrototypy(strona, rodziny) {
  if (rodziny.length === 0) {
    return {
      ...strona,
      tresc: `  <h1>Prototypy okien</h1>
  <div class="dc-uwaga dc-uwaga--brak">
    <p><strong>Makiet okien nie ma w wyjściu witryny.</strong> Strona bierze je z dostawy
    warstwy wizualnej (<code>design/05-okna/</code>) przy składaniu; skoro ich tu nie ma,
    dostawa nie była dostępna w chwili składania. Nie ma tu podglądów zastępczych: rysunek
    udający okno produktu byłby gorszy niż brak rysunku.</p>
  </div>`,
    };
  }

  const liczba = rodziny.reduce((suma, r) => suma + r.okna.length, 0);
  const pierwsze = rodziny[0].okna[0];

  const przyciski = rodziny
    .map(
      (r, i) =>
        `      <button type="button" class="dc-rodzina" data-rodzina="${tekst(r.klucz)}" role="tab" aria-selected="${i === 0 ? 'true' : 'false'}">${tekst(r.nazwa)} <span class="dane">${r.okna.length}</span></button>`,
    )
    .join('\n');

  const wykazy = rodziny
    .map(
      (r, i) => `    <ul class="dc-okna" data-rodzina="${tekst(r.klucz)}"${i === 0 ? '' : ' hidden'}>
${r.okna
  .map(
    (o) => `      <li><a class="dc-okno" href="${tekst(o.adres)}" data-tytul="${tekst(o.tytul)}"${o === pierwsze ? ' aria-current="true"' : ''}>${tekst(o.tytul)}<small>${tekst(o.plik)}</small></a></li>`,
  )
  .join('\n')}
    </ul>`,
    )
    .join('\n');

  const mobilne = rodziny
    .flatMap((r) => r.okna)
    .filter((o) => o.wlasna)
    .map(
      (o) => `    <figure class="dc-zrzut">
      <div class="dc-scena dc-scena--wysoka" data-szer="1100" data-wys="1500" style="--szer:1100px;--wys:1500px">
        <iframe class="dc-ramka" src="${tekst(o.adres)}" title="Prototyp okna: ${tekst(o.tytul)}" loading="lazy"></iframe>
      </div>
      <figcaption>${tekst(o.tytul)} — <a href="${tekst(o.adres)}">otwórz w pełnym oknie</a></figcaption>
    </figure>`,
    )
    .join('\n');

  return {
    ...strona,
    tresc: `  <style>${STYL}</style>
  <h1>Prototypy okien</h1>
  <p class="wiodacy">${liczba} okien produktu — <strong>działających</strong>, nie obrazków.
  Zakładki się przełączają, panele rozwijają, przełączniki działają. To makiety warstwy
  wizualnej, nie zrzuty z uruchomionego programu: pokazują kształt okna i jego zachowanie,
  a nie dane z czyjejś pracy.</p>

  <div class="dc-prototyp">
    <div class="dc-prototyp__wybor">
      <p class="dc-scena__podpis">Rodzina okien</p>
      <div class="dc-rodziny" role="tablist">
${przyciski}
      </div>
${wykazy}
    </div>
    <div class="dc-scena__wrap">
      <div class="dc-scena__pasek">
        <span class="dc-scena__podpis" id="podpis-prototypu">${tekst(pierwsze.tytul)}</span>
        <a class="dc-scena__pelne tlink" id="pelne-okno" href="${tekst(pierwsze.adres)}">Otwórz w pełnym oknie ↗</a>
      </div>
      <div class="dc-scena" data-szer="1440" data-wys="900" style="--szer:1440px;--wys:900px">
        <iframe class="dc-ramka" id="ramka-prototypu" src="${tekst(pierwsze.adres)}" title="Prototyp okna: ${tekst(pierwsze.tytul)}"></iframe>
      </div>
    </div>
  </div>
  <script>${SKRYPT}</script>

  <h2>Postać mobilna</h2>
  <p>Platforma ma opracowaną postać na wąski ekran oraz <strong>ekran zawsze widoczny</strong>.
  Poniżej stoją te dwie makiety — działające, jak pozostałe. Ekran telefonu widać w środku każdej
  z nich: dostawa opisuje postać mobilną jako funkcję globalną, więc makieta pokazuje ją wraz
  z opisem, a nie jako samo okno telefonu.</p>
  <p><strong>Instalki na telefon dziś nie ma</strong> — na stronie <a href="pobierz.html">Pobierz</a>
  stoją wyłącznie wydania na Windows i Linuksa. Pokazujemy tu kształt, nie zapowiedź terminu.</p>
  <div class="dc-zrzuty">
${mobilne}
  </div>

  <h2>Skąd te okna</h2>
  <p>Wszystkie pochodzą z dostarczonej warstwy wizualnej i wchodzą do witryny bez przeróbek —
  wraz z własnymi stylami, krojami i skryptami. Nazwa każdego okna czytana jest z samej makiety,
  więc wykaz obok nie może się rozjechać z tym, co się otwiera.</p>`,
  };
}
