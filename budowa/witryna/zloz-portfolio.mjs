// Złożenie strony produktu DLA WITRYNY GRUPY — `danaco-group.pl/danaco-console.html`.
//
// PO CO DRUGIE POLECENIE, SKORO JEST `zloz.mjs`. Bo to są dwie różne witryny
// i dwa różne języki wizualne, a nie dwa warianty jednej.
//
//   `zloz.mjs`          → witryna produktu na kanale wydań (pobierz.danaco-group.pl):
//                         żetony `design/`, monochromatyczna precyzja, atramentowy
//                         pasek, strona „Pobierz" z sumami SHA-256 i przyciskami.
//   `zloz-portfolio.mjs`→ JEDNA strona w portfolio grupy kapitałowej: arkusz
//                         tamtej witryny (jasny, złote akcenty, kursywa
//                         szeryfowa), jej nagłówek i stopka wzięte dosłownie,
//                         BEZ plików do pobrania — z jednym przejściem do kanału.
//
// Kanon Console rządzi produktem i tym, co widać w RAMKACH prototypów. Nie rządzi
// witryną grupy: tam obowiązuje jej własna identyfikacja, bo strona ma być
// częścią tamtej witryny, a nie wyspą we własnym stylu.
//
// JEDNA STRONA, DWA WEJŚCIA. Produkt wchodzi do menu dwa razy — z PORTFOLIO
// („co rozwijamy") i z DANACO SHARE („jak to wziąć") — ale strona jest jedna,
// z zakotwiczeniami `#dorobek`, `#platforma`, `#pobierz`. Dwie strony o tym samym
// produkcie to dwie prawdy, które rozjadą się przy pierwszej poprawce.
//
// TREŚĆ POCHODZI Z OPRACOWAŃ. Nazwy środowisk, tryby pracy i przeznaczenie —
// z `docs/architektura/koncepcja-platformy.md`, rozdz. 9. Wykaz modułów — z nazw
// opracowań w `docs/moduly/`. Metryka — z `design/KANON.md`, rozdz. 0. Nic tu nie
// jest dopisane z głowy wykonawcy: kanon tego zabrania, a strona wizerunkowa
// grupy kapitałowej jest ostatnim miejscem, w którym wolno coś zmyślić.
import { mkdir, writeFile, readFile, readdir, rm, cp } from 'node:fs/promises';
import { dirname, join } from 'node:path';
import { fileURLToPath } from 'node:url';

import { katalogPrototypow } from './tresc/prototypy.mjs';

const KATALOG = dirname(fileURLToPath(import.meta.url));
const DOSTAWA = join(KATALOG, '..', '..', 'design');
const OPRACOWANIA = join(KATALOG, '..', '..', 'docs');
const PACZKA = join(KATALOG, '..', '..', '..', 'DanacoConsole.git', 'DO-WGRANIA-DANACO-GROUP');
const KANAL = 'https://pobierz.danaco-group.pl';

function tekst(wartosc) {
  return String(wartosc)
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;');
}

// Cztery środowiska — przepisane z tabeli rozdz. 9 opracowania koncepcji platformy.
// Trzecia kolumna („Boczna nawigacja") pominięta: na stronie wizerunkowej nie
// niesie treści. Nic nie dodano.
const SRODOWISKA = [
  { nazwa: 'TalkIn', tryb: 'Wiedza, komunikacja i praca z treścią', po_co: 'Rozmowy z AI, praca z dokumentami, analizy, tłumaczenia, badania, raporty, współpraca wielu modeli.' },
  { nazwa: 'WorkSpace', tryb: 'Produktywność, organizacja i realizacja projektów', po_co: 'Zarządzanie projektami, automatyzacja, prowadzenie procesów, tworzenie aplikacji biznesowych, projektowanie, praca operacyjna.' },
  { nazwa: 'CodeStudio', tryb: 'Programowanie', po_co: 'Tworzenie oprogramowania, praca z kodem, debugowanie, wykorzystanie terminali, budowanie aplikacji, procesy developerskie.' },
  { nazwa: 'MultitaskingAI', tryb: 'Orkiestracja autonomicznej pracy ciągłej', po_co: 'Zespół modeli i agentów w rolach; po spięciu z automatyką — pełna pętla pracy ciągłej z własnymi akcjami i harmonogramem.' },
];

// Postaci produktu — z rozstrzygnięcia modelu wdrożenia (CLAUDE.md, rozdz. 2)
// oraz z wykazu wydań. Bez rozmiarów i bez sum: to jest strona wizerunkowa,
// a sumy i pliki mają jedno miejsce — kanał wydań.
const POSTACI = [
  ['Hybryda · Windows x64', 'Okno na stanowisko z procesorem Intel albo AMD. Rdzeń pracuje na serwerze wdrożenia.'],
  ['Hybryda · Windows ARM64', 'To samo okno dla maszyn z Windows na procesorze ARM.'],
  ['Hybryda · Linux', 'Dwie postacie pliku: samodzielny AppImage oraz pakiet dla menedżera pakietów.'],
  ['Natywna pełna · Linux', 'Całość na jednym urządzeniu — okno, powłoka i rdzeń. Serwer nie jest potrzebny.'],
  ['Natywna pełna · Windows', 'Odpowiednik postaci natywnej dla Windows, w jednej instalce.'],
];

/** Wykaz modułów z nazw opracowań w `docs/moduly/` — nie z listy wpisanej ręcznie.
 *  Piętnaście modułów jest w opracowaniu koncepcji liczbą wiążącą, więc liczba
 *  na stronie ma pochodzić z pomiaru katalogu, a nie z przepisania. */
async function moduly() {
  try {
    const pliki = await readdir(join(OPRACOWANIA, 'moduly'));
    return pliki
      .filter((p) => p.endsWith('.md'))
      .map((p) => p.replace(/\.md$/, ''))
      .map((n) => n.charAt(0).toUpperCase() + n.slice(1))
      .sort();
  } catch {
    return [];
  }
}

// Jedyny własny styl na tej stronie — i tylko mechanika ramek. Arkusz witryny
// grupy nie zna pojęcia „prototyp w ramce przeskalowanej do szerokości", bo
// wcześniej nie było na niej prototypów. Wszystko, co da się wyrazić jej klasami
// (`section`, `container`, `area`, `tag`, `prose`, `btn`), wyrażone jest nimi —
// tu zostaje sama scena, ramka i wykaz okien. Osobnego pliku CSS NIE dokładamy:
// jeden arkusz witryny zostaje jednym arkuszem, a to jest kilkanaście reguł.
const STYL = `<style>
.dc-prototyp{display:grid;grid-template-columns:17rem 1fr;gap:var(--space-lg,24px);align-items:start}
.dc-rodziny{display:flex;flex-wrap:wrap;gap:8px;margin-bottom:16px}
.dc-rodzina{font-family:var(--ff-body);font-size:.82rem;letter-spacing:.02em;min-height:34px;padding:0 14px;cursor:pointer;background:transparent;color:var(--text-muted);border:1px solid var(--border-default);border-radius:var(--radius-pill)}
.dc-rodzina[aria-selected="true"]{background:var(--bg-inverse);color:var(--bg-surface);border-color:var(--bg-inverse);font-weight:600}
.dc-okna{list-style:none;margin:0;padding:0;border:1px solid var(--border-default);border-radius:var(--radius-lg);overflow:hidden;background:var(--bg-surface)}
.dc-okna li+li{border-top:1px solid var(--border-default)}
.dc-okno{display:block;padding:10px 14px;color:var(--text-body);text-decoration:none;font-size:.86rem;line-height:1.35}
.dc-okno:hover{background:var(--bg-subtle)}
/* Kursywa szeryfowa (Fraunces) w wyróżnieniu części zdania — tym samym wzorcem,
   co na żywej witrynie grupy: nigdy w całym akapicie, zawsze w akcencie. */
.dc-akcent em{font-family:var(--ff-editorial);font-style:italic;font-weight:400;color:var(--gold)}
.dc-okno small{display:block;color:var(--text-muted);font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;font-size:.74rem}
.dc-okno[aria-current="true"]{background:var(--bg-subtle);font-weight:600;box-shadow:inset 3px 0 0 var(--brand-yellow)}
.dc-scena{position:relative;overflow:hidden;border:1px solid var(--border-default);border-radius:var(--radius-lg);height:600px;background:var(--bg-surface)}
.dc-ramka{position:absolute;top:0;left:0;width:1440px;height:900px;border:0;transform-origin:0 0;display:block}
.dc-pasek{display:flex;align-items:center;gap:12px;flex-wrap:wrap;margin-bottom:8px;font-family:var(--ff-display);font-size:.74rem;letter-spacing:.08em;text-transform:uppercase;color:var(--text-muted)}
.dc-pasek a{margin-left:auto;text-transform:none;letter-spacing:0}
/* Makieta postaci mobilnej i ekranu zawsze widocznego NIE JEST dokumentem
   o szerokości telefonu — to strona opisowa (do 1100 px), która pokazuje w środku
   ekran telefonu wraz z opisem funkcji globalnej. Wciśnięta w ramkę o proporcjach
   telefonu rozjeżdżała się, bo ma własny próg 1100 px. Dlatego idzie w szerokości
   własnej i jest wyższa niż szersza. */
.dc-wysoka{height:860px}
.dc-wysoka .dc-ramka{width:1100px;height:1500px}
@media (max-width:1100px){
  .dc-prototyp{grid-template-columns:1fr}
  .dc-scena, .dc-wysoka{height:58vh;min-height:300px}
}
</style>`;

const SKRYPT = `<script>
(function(){
  var ramka=document.getElementById('dc-ramka'), podpis=document.getElementById('dc-podpis'), pelne=document.getElementById('dc-pelne');
  function skaluj(){
    var sceny=document.querySelectorAll('.dc-scena');
    for(var i=0;i<sceny.length;i+=1){
      var scena=sceny[i], r=scena.querySelector('.dc-ramka');
      if(!r) continue;
      var szer=parseInt(scena.getAttribute('data-szer'),10)||1440;
      var wys=parseInt(scena.getAttribute('data-wys'),10)||900;
      var s=scena.clientWidth/szer;
      if(!isFinite(s)||s<=0) continue;
      r.style.transform='scale('+s+')';
      scena.style.height=Math.round(wys*s)+'px';
    }
  }
  if(ramka&&podpis&&pelne){
    document.addEventListener('click',function(e){
      var a=e.target.closest?e.target.closest('.dc-okno'):null;
      if(!a) return;
      if(e.metaKey||e.ctrlKey||e.shiftKey||e.button!==0) return;
      e.preventDefault();
      ramka.setAttribute('src',a.getAttribute('href'));
      var nazwa=a.getAttribute('data-tytul')||a.textContent.trim();
      ramka.setAttribute('title','Prototyp okna: '+nazwa);
      podpis.textContent=nazwa;
      pelne.setAttribute('href',a.getAttribute('href'));
      var w=document.querySelectorAll('.dc-okno');
      for(var i=0;i<w.length;i+=1){w[i].removeAttribute('aria-current');}
      a.setAttribute('aria-current','true');
      if(window.matchMedia('(max-width:960px)').matches){ramka.scrollIntoView({block:'start',behavior:'smooth'});}
      window.setTimeout(skaluj,0);
    });
  }
  var przyciski=document.querySelectorAll('.dc-rodzina');
  function pokaz(k){
    for(var i=0;i<przyciski.length;i+=1){przyciski[i].setAttribute('aria-selected',przyciski[i].getAttribute('data-rodzina')===k?'true':'false');}
    var listy=document.querySelectorAll('.dc-okna');
    for(var j=0;j<listy.length;j+=1){listy[j].hidden=listy[j].getAttribute('data-rodzina')!==k;}
    window.setTimeout(skaluj,0);
  }
  for(var k=0;k<przyciski.length;k+=1){przyciski[k].addEventListener('click',function(){pokaz(this.getAttribute('data-rodzina'));});}
  if(przyciski.length){pokaz(przyciski[0].getAttribute('data-rodzina'));}
  skaluj();
  window.addEventListener('resize',skaluj);
})();
</script>`;

function sekcjaPrototypow(rodziny) {
  if (rodziny.length === 0) {
    return `      <p class="prose">Prototypów okien nie dołączono do tej strony — katalog makiet nie był dostępny
      przy jej składaniu. Nie ma tu rysunków zastępczych: obrazek udający okno produktu byłby
      gorszy niż jego brak.</p>`;
  }
  const liczba = rodziny.reduce((s, r) => s + r.okna.length, 0);
  const pierwsze = rodziny[0].okna[0];
  const przyciski = rodziny
    .map(
      (r, i) =>
        `            <button type="button" class="dc-rodzina" data-rodzina="${tekst(r.klucz)}" aria-selected="${i === 0 ? 'true' : 'false'}">${tekst(r.nazwa)} (${r.okna.length})</button>`,
    )
    .join('\n');
  const listy = rodziny
    .map(
      (r, i) => `          <ul class="dc-okna" data-rodzina="${tekst(r.klucz)}"${i === 0 ? '' : ' hidden'}>
${r.okna
  .map(
    (o) =>
      `            <li><a class="dc-okno" href="${tekst(o.adres)}" data-tytul="${tekst(o.tytul)}"${o === pierwsze ? ' aria-current="true"' : ''}>${tekst(o.tytul)}<small>${tekst(o.plik)}</small></a></li>`,
  )
  .join('\n')}
          </ul>`,
    )
    .join('\n');
  return `      <p class="prose" style="max-width:64ch;">Poniżej stoi ${liczba} okien platformy —
      <strong>działających</strong>, nie zrzutów ekranu. Zakładki się przełączają, panele rozwijają,
      przełączniki działają. To prototypy warstwy wizualnej produktu: pokazują kształt okna i jego
      zachowanie, nie dane z czyjejś pracy.</p>
      <div class="dc-prototyp reveal" style="margin-top:var(--space-lg);">
        <div>
          <div class="dc-rodziny">
${przyciski}
          </div>
${listy}
        </div>
        <div>
          <div class="dc-pasek"><span id="dc-podpis">${tekst(pierwsze.tytul)}</span><a id="dc-pelne" class="tlink" href="${tekst(pierwsze.adres)}">Otwórz w pełnym oknie →</a></div>
          <div class="dc-scena" data-szer="1440" data-wys="900">
            <iframe class="dc-ramka" id="dc-ramka" src="${tekst(pierwsze.adres)}" title="Prototyp okna: ${tekst(pierwsze.tytul)}"></iframe>
          </div>
        </div>
      </div>`;
}

function sekcjaMobilna(rodziny) {
  const okna = rodziny.flatMap((r) => r.okna).filter((o) => o.wlasna);
  if (okna.length === 0) return '';
  return okna
    .map(
      (o) => `        <figure style="margin:0 0 var(--space-lg,24px);">
          <div class="dc-pasek"><span>${tekst(o.tytul)}</span><a class="tlink" href="${tekst(o.adres)}">Otwórz w pełnym oknie →</a></div>
          <div class="dc-scena dc-wysoka" data-szer="1100" data-wys="1500">
            <iframe class="dc-ramka" src="${tekst(o.adres)}" title="Prototyp okna: ${tekst(o.tytul)}" loading="lazy"></iframe>
          </div>
        </figure>`,
    )
    .join('\n');
}

const rodziny = await katalogPrototypow(DOSTAWA);
const wykazModulow = await moduly();

const TRESC = `
  <section class="page-head">
    <div class="container">
      <span class="kicker">Portfolio rozwoju · technologia</span>
      <h1 class="dc-akcent">Danaco <em>Console.</em></h1>
      <p class="lead">Platforma AI Workspace OS — wielośrodowiskowy system operacyjny dla sztucznej
      inteligencji, integrujący komunikację, zarządzanie wiedzą, tworzenie treści, projektowanie,
      automatyzacje procesów oraz rozwój oprogramowania.</p>
      <a class="tlink hero__link" href="#dorobek">Co rozwijamy w tym produkcie <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M13 6l6 6-6 6"/></svg></a>

      <div class="page-head__aside reveal">
        <div class="ttl">Metryka</div>
        <div class="row"><span class="k">Producent</span><span class="v">Danaco Holding Group</span></div>
        <div class="row"><span class="k">Twórca</span><span class="v">Dariusz Naharnowicz</span></div>
        <div class="row"><span class="k">Środowiska</span><span class="v">${SRODOWISKA.length}</span></div>
        <div class="row"><span class="k">Moduły</span><span class="v">${wykazModulow.length}</span></div>
        <div class="row"><span class="k">Status</span><span class="v">deweloperski</span></div>
      </div>
    </div>
  </section>

  <section class="section section--paper" id="dorobek">
    <div class="container">
      <div class="section-head reveal">
        <div class="secno"><b>Co rozwijamy</b></div>
        <h2 class="dc-akcent" style="margin-top:var(--space-lg);">Cztery środowiska pracy, ${wykazModulow.length} modułów, <em>jedno stanowisko.</em></h2>
      </div>
      <div class="split" style="margin-top:var(--space-lg);">
        <div class="prose reveal">
          <p>Danaco Console rozwijamy jako stanowisko pracy z modelami AI dla zawodowego Operatora.
          Punktem wyjścia nie jest okno rozmowy, do którego dokłada się kolejne funkcje, ale
          <strong>zadanie</strong> — a dla zadania dobierana jest osobna przestrzeń robocza. Wraz
          z zadaniem zmienia się układ okien, dostępne narzędzia, historia sesji, pamięć kontekstowa
          i sposób działania modelu.</p>
          <p>Najwyższym poziomem organizacji pracy jest <strong>środowisko</strong>: definiuje układ
          interfejsu, dostępne moduły, typ pracy i model nawigacji. Wewnątrz środowisk pracują
          <strong>moduły</strong> — wyspecjalizowane obszary robocze z własnym zestawem okien
          operacyjnych i własnym przepływem pracy.</p>
        </div>
        <div class="prose reveal">
          <p>Ponad środowiskami i modułami działa warstwa <strong>funkcji globalnych</strong>:
          postać mobilna oraz ekran zawsze widoczny. Nad wszystkim stoi zasada, którą produkt
          egzekwuje w kodzie: <strong>odmowa nazywa brak</strong> — jeżeli czegoś nie ma, platforma
          mówi to wprost i podaje, czego brakuje, zamiast zwracać wynik pozorny.</p>
          <p>Warstwa wizualna produktu została opracowana w wersji 2.0 (11.08.2026) w języku
          <em>monochromatycznej precyzji</em>: odcienie bieli w motywie jasnym, czerni w ciemnym,
          jeden chłodny sygnał, dane krojem mono, zero dekoracji bez funkcji.</p>
        </div>
      </div>

      <div class="section-head reveal" style="margin-top:var(--space-xl,48px);">
        <div class="secno"><b>Środowiska</b></div>
      </div>
      <div class="grid grid-2" style="margin-top:var(--space-lg);">
${SRODOWISKA.map(
  (s) => `        <article class="area reveal">
          <div class="area__top"><svg class="chev" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 6l6 6-6 6M13 6l6 6-6 6"/></svg></div>
          <span class="tag">${tekst(s.tryb)}</span>
          <h3>${tekst(s.nazwa)}</h3>
          <p>${tekst(s.po_co)}</p>
        </article>`,
).join('\n')}
      </div>

      <div class="section-head reveal" style="margin-top:var(--space-xl,48px);">
        <div class="secno"><b>Moduły platformy</b></div>
      </div>
      <div class="flow" style="margin-top:var(--space-lg);">
${wykazModulow.map((m) => `        <div class="flow__step"><span class="n">${tekst(m)}</span></div>`).join('\n')}
      </div>
      <p class="figure__cap">${wykazModulow.length} modułów; każdy dostępny w wybranych środowiskach.</p>
    </div>
  </section>

  <section class="section section--paper" id="prototypy">
    <div class="container">
      <div class="section-head reveal">
        <div class="secno"><b>Prototypy okien</b></div>
        <h2 class="dc-akcent" style="margin-top:var(--space-lg);">Produkt w pracy, <em>a nie opis produktu.</em></h2>
      </div>
${sekcjaPrototypow(rodziny)}
    </div>
  </section>

  <section class="section section--paper" id="platforma">
    <div class="container">
      <div class="section-head reveal">
        <div class="secno"><b>Jak pracuje platforma</b></div>
        <h2 class="dc-akcent" style="margin-top:var(--space-lg);">Serce na serwerze, <em>cienka instalka u klienta.</em></h2>
      </div>
      <div class="split" style="margin-top:var(--space-lg);">
        <div class="prose reveal">
          <p>Platforma jest <strong>hybrydą</strong>. Całość — rdzeń, baza, magazyn zasobów i arsenał
          narzędzi — stoi na serwerze wdrożenia. U klienta zakłada się <strong>cienką instalkę</strong>:
          natywne okno, które po uruchomieniu pokazuje ekran logowania, a po wejściu całą platformę.
          Operator nie zakłada u siebie żadnego programu z arsenału.</p>
          <p>Osobno istnieje <strong>wersja całkowicie natywna</strong> — dla Windows i dla Linuksa —
          w której rdzeń i arsenał stoją na urządzeniu klienta. Jest większa, ale nie potrzebuje
          serwera po drugiej stronie.</p>
        </div>
        <div class="prose reveal">
          <p>Praca Operatora nie opuszcza jego infrastruktury: rozmowy, pliki i historia leżą tam,
          gdzie powstały. Producent nie pośredniczy w treści pracy.</p>
          <p><strong>Postać mobilna</strong> i <strong>ekran zawsze widoczny</strong> są opracowane
          jako funkcje globalne platformy — poniżej stoją ich działające prototypy. Instalki na
          telefon dziś nie ma i nie zapowiadamy tu terminu.</p>
        </div>
      </div>
      <div style="margin-top:var(--space-lg);">
${sekcjaMobilna(rodziny)}
      </div>
    </div>
  </section>

  <section class="section section--paper" id="pobierz">
    <div class="container">
      <div class="section-head reveal">
        <div class="secno"><b>Dostęp do platformy</b></div>
        <h2 class="dc-akcent" style="margin-top:var(--space-lg);">Pięć postaci produktu <em>oraz pakiet serwera.</em></h2>
      </div>
      <div class="prose reveal" style="margin-top:var(--space-lg);max-width:64ch;">
        <p>Wydania platformy leżą na osobnym kanale — <strong>nie na tej stronie</strong>. Tu jest
        wykaz postaci, w których produkt się wydaje; pliki, ich rozmiary i sumy kontrolne SHA-256
        mają jedno miejsce i jest nim kanał wydań.</p>
      </div>
      <div class="grid grid-2" style="margin-top:var(--space-lg);">
${POSTACI.map(
  ([nazwa, opis]) => `        <article class="area reveal">
          <div class="area__top"><svg class="chev" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.6" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M6 6l6 6-6 6M13 6l6 6-6 6"/></svg></div>
          <h3>${tekst(nazwa)}</h3>
          <p>${tekst(opis)}</p>
        </article>`,
).join('\n')}
      </div>
      <div class="prose reveal" style="margin-top:var(--space-lg);max-width:64ch;">
        <p>Osobno wydawany jest <strong>pakiet serwera wdrożenia</strong> — zakłada go administrator
        na maszynie serwerowej i dopiero wtedy postacie hybrydowe mają z czym rozmawiać. Klient
        pracujący sam na jednym komputerze bierze postać natywną pełną i serwera nie potrzebuje.</p>
        <p><strong>Pobieranie jest chronione hasłem i na razie zastrzeżone.</strong> Po przejściu na
        kanał wydań przeglądarka poprosi o nazwę użytkownika i hasło — nie jest to usterka i nie ma
        tam formularza rejestracji. Poświadczenie wydaje producent:
        <a href="mailto:support@danaco-group.pl">support@danaco-group.pl</a>.</p>
      </div>
      <div class="hero__actions reveal" style="margin-top:var(--space-lg);">
        <a class="btn btn--primary" href="${KANAL}/pobierz.html" rel="noopener">Przejdź do kanału wydań <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12h14M13 6l6 6-6 6"/></svg></a>
      </div>
    </div>
  </section>
${SKRYPT}`;

const szablon = await readFile(join(KATALOG, 'portfolio', 'szkielet-grupy.html'), 'utf8');
const strona = szablon
  .replaceAll('{{TYTUL}}', 'Danaco Console — platforma AI Workspace OS | Danaco Group')
  .replaceAll(
    '{{OPIS}}',
    'Danaco Console — platforma AI Workspace OS rozwijana w Danaco Group: cztery środowiska pracy, moduły, postać mobilna oraz działające prototypy okien.',
  )
  .replace('{{STYL}}', STYL)
  .replace('{{TRESC}}', TRESC);

await rm(PACZKA, { recursive: true, force: true });
await mkdir(join(PACZKA, 'prototypy'), { recursive: true });
await writeFile(join(PACZKA, 'danaco-console.html'), strona, 'utf8');

for (const czlon of ['05-okna', '_samodzielne', 'zasoby']) {
  await cp(join(DOSTAWA, czlon), join(PACZKA, 'prototypy', czlon), { recursive: true });
}

// Fragmenty menu — osobno, bo wkleja się je w DWA różne miejsca i w KAŻDYM
// pliku tamtej witryny. Podanie ich jako jednego bloku kończyłoby się wklejeniem
// obu w jedno rozwinięcie.
await writeFile(
  join(PACZKA, 'FRAGMENT-MENU-PORTFOLIO.html'),
  `<!-- POZYCJA MENU: rozwinięcie „Portfolio" -->
<!-- Wkleić w KAŻDYM pliku .html witryny, w bloku <div class="nav__dd-list"> należącym -->
<!-- do pozycji „Portfolio" (ta, której nav__link prowadzi do obszary.html).            -->
<!-- Miejsce: BEZPOŚREDNIO PRZED wierszem z „Nasze marki hotelowe".                     -->
              <a class="nav__dd-link" href="danaco-console.html" role="menuitem">Danaco Console</a>
`,
  'utf8',
);
await writeFile(
  join(PACZKA, 'FRAGMENT-MENU-SHARE.html'),
  `<!-- POZYCJA MENU: rozwinięcie „Danaco Share" -->
<!-- Wkleić w KAŻDYM pliku .html witryny, w bloku <div class="nav__dd-list"> należącym  -->
<!-- do pozycji „Danaco Share" (ta, której nav__link prowadzi do aktualnosci.html).     -->
<!-- Miejsce: BEZPOŚREDNIO PRZED wierszem z „Repozytorium dokumentów".                  -->
<!-- Odnośnik prowadzi do sekcji #platforma, bo z tego menu wchodzi się po tym,         -->
<!-- jak platforma pracuje i jak ją wziąć — nie po dorobku.                             -->
              <a class="nav__dd-link" href="danaco-console.html#platforma" role="menuitem">Danaco Console</a>
`,
  'utf8',
);

await writeFile(
  join(PACZKA, 'FRAGMENT-MENU-TOP.html'),
  `<!-- POZYCJA MENU: pasek górny, obok O NAS · DZIAŁALNOŚĆ · PORTFOLIO · DANACO SHARE -->
<!-- Wkleić w KAŻDYM pliku .html witryny, wewnątrz <nav class="nav__menu">.            -->
<!-- Miejsce: BEZPOŚREDNIO PRZED wierszem z pozycją „Kontakt" — kontakt zwyczajowo     -->
<!-- zamyka pasek, więc Pobierz staje przed nim.                                        -->
<!--                                                                                    -->
<!-- Bez rozwinięcia (\`nav__dd\`): cztery sąsiednie pozycje mają podmenu, ta nie ma     -->
<!-- czego w nim trzymać — prowadzi wprost do plików.                                    -->
<!-- Bez \`target="_blank"\`: subdomena to nadal witryna grupy, a otwieranie nowej karty -->
<!-- przy przejściu w obrębie własnych stron jest natrętne.                              -->
      <a class="nav__link" href="https://pobierz.danaco-group.pl/" rel="noopener">Pobierz</a>
`,
  'utf8',
);

const liczbaOkien = rodziny.reduce((s, r) => s + r.okna.length, 0);

await writeFile(
  join(PACZKA, 'JAK-WGRAC.md'),
  `# Jak wgrać stronę Danaco Console na danaco-group.pl

Paczka jest kompletem do przeniesienia. Dostępu SSH do hosta \`danaco-group.pl\`
(137.74.41.149) wykonawca nie ma — dlatego to jest paczka, a nie wdrożenie.

Złożone: ${new Date().toISOString().slice(0, 10)}.

---

## 1. Co tu leży

| Element | Co to jest | Gdzie kopiować |
|---|---|---|
| \`danaco-console.html\` | strona o produkcie — dorobek, prototypy, jak platforma pracuje, przejście do pobierania | korzeń witryny, obok \`obszary.html\` |
| \`prototypy/\` | ${liczbaOkien} działających prototypów okien wraz z ich zasobami (styl, kroje, ikony, skrypty) | korzeń witryny, jako \`prototypy/\` |
| \`FRAGMENT-MENU-PORTFOLIO.html\` | pozycja do wklejenia w rozwinięciu **Portfolio** | nie kopiować — wkleić treść |
| \`FRAGMENT-MENU-SHARE.html\` | pozycja do wklejenia w rozwinięciu **Danaco Share** | nie kopiować — wkleić treść |
| \`FRAGMENT-MENU-TOP.html\` | pozycja **Pobierz** w pasku górnym | nie kopiować — wkleić treść |

**Czego tu nie ma i nie ma być:** waszego arkusza \`assets/css/styles.css\`.
Strona **linkuje** go, nie nosi kopii — kopia rozjechałaby się przy pierwszej
waszej poprawce i \`danaco-console.html\` zostałaby jedyną podstroną w starym
wyglądzie. W paczce nie ma też katalogu \`assets/\`: strona nie dokłada ani jednego
własnego zasobu poza prototypami.

## 2. Kopiowanie

\`\`\`
danaco-console.html   →  /danaco-console.html
prototypy/            →  /prototypy/
\`\`\`

Układ katalogów wewnątrz \`prototypy/\` **musi zostać zachowany**. Makiety wołają
zasoby ścieżkami względnymi (\`../../zasoby/css/fundament.css\`), więc przeniesienie
samych plików HTML da okna bez stylów. Kopiuj cały katalog, nie jego zawartość
po kawałku.

## 3. Trzy pozycje menu — wklejane w KAŻDYM pliku .html

Menu jest wpisane w każdą stronę witryny osobno (nie jest dołączane), więc każdą
z trzech pozycji trzeba wkleić tyle razy, ile jest plików: \`index.html\`,
\`o-nas.html\`, \`obszary.html\`, \`inwestycje.html\`, \`aktualnosci.html\`,
\`historia.html\`, \`kariera.html\`, \`kontakt.html\`, \`dokumenty.html\` oraz
w podstronach katalogu \`dokumenty/\`.

| Fragment | Blok docelowy | Wstawić bezpośrednio przed wierszem |
|---|---|---|
| \`FRAGMENT-MENU-PORTFOLIO.html\` | \`<div class="nav__dd-list">\` pozycji **Portfolio** | \`… href="obszary.html#hotelarstwo" …>Nasze marki hotelowe</a>\` |
| \`FRAGMENT-MENU-SHARE.html\` | \`<div class="nav__dd-list">\` pozycji **Danaco Share** | \`… href="dokumenty.html" …>Repozytorium dokumentów</a>\` |
| \`FRAGMENT-MENU-TOP.html\` | \`<nav class="nav__menu">\` | \`<a class="nav__link" href="kontakt.html">Kontakt</a>\` |

Trzy wejścia, trzy różne cele — i dlatego są trzy, a nie jedno:

| Wejście | Po co | Adres |
|---|---|---|
| Portfolio → Danaco Console | dorobek: środowiska, moduły, prototypy | \`danaco-console.html\` |
| Danaco Share → Danaco Console | jak platforma pracuje i jak ją wziąć | \`danaco-console.html#platforma\` |
| Pobierz (pasek górny) | wprost do plików instalek | \`https://pobierz.danaco-group.pl/\` |

**Strona jest jedna, nie trzy.** Dwie strony o tym samym produkcie to dwie prawdy,
które rozjadą się przy pierwszej poprawce; stąd zakotwiczenia \`#dorobek\`,
\`#platforma\`, \`#pobierz\` w jednym pliku.

## 4. Skąd wzięty jest wygląd tej strony

**Kanonu witryny grupy nie ma w postaci opracowania.** Nie leży ani w \`design/\`,
ani w \`docs/\` — jedynym źródłem jest żywa witryna. Wszystko, co o nim przyjęto,
wyprowadzono z dwóch rzeczy:

- arkusza \`https://www.danaco-group.pl/assets/css/styles.css?v=20260602f\`
  (137 KB) — barwy marki, powierzchnie obu motywów, kroje (General Sans na
  nagłówki, Manrope na tekst, Fraunces na kursywę akcentu), żetony przycisków,
  progi widoku liczone po \`max-width\`;
- markupu \`https://www.danaco-group.pl/obszary.html\` — nagłówek i stopka wzięte
  **dosłownie** (wycina je polecenie \`portfolio/wytnij-szkielet.mjs\`, nie ręczne
  przepisanie), oraz klasy sekcji: \`section\`, \`container\`, \`section-head\`,
  \`secno\`, \`page-head\`, \`grid grid-2\`, \`area\`, \`tag\`, \`prose\`, \`split\`,
  \`flow\`, \`btn btn--primary\`, \`kicker\`, \`lead\`, \`tlink\`.

**Wersja arkusza jest tym, wobec czego strona była składana: \`v=20260602f\`.**
Gdy zmienicie tę wersję, stronę należy przejrzeć — nie zakładać, że nadal pasuje.
Gdyby powstała księga marki grupy, ona będzie źródłem, nie ten odczyt z CSS.

Własnych reguł stylu strona niesie kilkanaście, w jednym bloku \`<style>\` w jej
nagłówku, i tylko na rzeczy, których wasz arkusz nie zna: wykaz okien, scena
z prototypem i skalowanie ramki. Osobnego pliku CSS nie dokłada — jeden arkusz
witryny zostaje jednym arkuszem.

## 5. Granica dwóch kanonów

- **Strona** — kanon danaco-group: wasz arkusz, wasze kroje, złote akcenty,
  kursywa szeryfowa w wyróżnieniu części zdania.
- **Wnętrze ramek z prototypami** — kanon Danaco Console (monochromatyczna
  precyzja, Space Grotesk i IBM Plex, jeden chłodny sygnał). Tak ma być: ramka
  pokazuje produkt, nie stronę, i niesie własny \`fundament.css\` oraz
  \`prototyp.css\`, więc dwa drzewa stylów nie mieszają się ze sobą.

Poza wnętrzem ramek na stronie nie ma ani jednego żetonu \`--dn-*\`, ani jednego
kroju z dostawy Console, ani jednej barwy z tamtego kanonu.

## 6. Sprawdzenie po wgraniu

\`\`\`bash
curl -sI https://www.danaco-group.pl/danaco-console.html           # 200
curl -sI https://www.danaco-group.pl/prototypy/05-okna/moduly/studio.html   # 200
curl -sI https://www.danaco-group.pl/prototypy/zasoby/prototyp.js           # 200
\`\`\`

W przeglądarce sprawdzić trzy rzeczy, których \`curl\` nie pokaże:

1. **Prototyp w ramce rysuje się ze stylami** — jeżeli okno jest białe i bez
   układu, katalog \`prototypy/zasoby/\` nie doszedł albo doszedł w innym miejscu.
2. **Motyw ciemny** — wasz arkusz niesie oba motywy; strona ma znosić oba.
   Przełączenie motywu systemu nie może zostawić białych plam.
3. **Trzy pozycje menu** widać na każdej stronie, nie tylko na nowej.

### Jedna rzecz zmierzona po waszej stronie, nie po naszej

W arkuszu \`styles.css\` dwie reguły globalne wiążą barwę na stałe:
\`strong { color: var(--ink) }\` oraz \`h3, h4, h5 { color: var(--ink) }\`, a zmienna
\`--ink\` (#18233B) **nie przełącza się z motywem**. W motywie ciemnym daje to
granat na granacie — dotyczy to również waszych istniejących stron (nagłówki kart
\`.area\` na \`obszary.html\`), nie tylko tej nowej. Nie ruszaliśmy tego: nadpisanie
waszej reguły globalnie byłoby wejściem w wasz arkusz. Na stronie kanału wydań
(\`pobierz.danaco-group.pl\`) poprawka jest zawężona do własnych komponentów
i użyto w niej **waszego** żetonu \`--text-strong\`, który motyw przełącza.
Gdybyście chcieli to naprawić u siebie, wystarczy jedna zmiana: \`--ink\` na
\`var(--text-strong)\` w tych dwóch regułach.

## 7. Treść — skąd pochodzi

- Nazwy środowisk, tryby pracy i przeznaczenie: \`docs/architektura/koncepcja-platformy.md\`, rozdz. 9.
- Wykaz modułów: nazwy opracowań w \`docs/moduly/\` (liczba czytana z katalogu, nie przepisana).
- Metryka (producent, twórca, kontakt, status **deweloperski**, warstwa wizualna v2.0 z 11.08.2026): \`design/KANON.md\`, rozdz. 0.
- Model wdrożenia (hybryda, wersja natywna, pakiet serwera): rozstrzygnięcie Właściciela.

\`design/KANON.md\` obowiązuje w treści: **zakaz wymyślania modułów, funkcji, nazw,
metryk i osób.** Na tej stronie nie ma liczb „z rynku" ani zapowiedzi terminów.
Postać mobilna jest pokazana jako opracowana funkcja globalna — z jawnym zdaniem,
że **instalki na telefon nie ma**.
`,
  'utf8',
);

console.log(`strona portfolio złożona: ${join(PACZKA, 'danaco-console.html')}`);
console.log(`prototypów osadzonych: ${liczbaOkien} (${rodziny.map((r) => `${r.nazwa}: ${r.okna.length}`).join(', ')})`);
console.log(`modułów z opracowań: ${wykazModulow.length}`);
