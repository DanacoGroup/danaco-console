// Strona „Pobierz" — powstaje z wykazu wydań, nie z napisanej treści.
//
// Pusty wykaz nie jest błędem. Strona bez wydań pokazuje zdanie „nie ma jeszcze
// żadnego wydania" — nie przycisk, który prowadzi donikąd, i nie „wkrótce"
// udające zapowiedź. Przycisk bez pliku byłby obietnicą bez pokrycia.
//
// Milczenie o wydaniu, które istnieje, jest tą samą nieprawdą w drugą stronę.
// Bywa, że pakiety leżą zbudowane, a kanał pobierania (DNS, serwer, certyfikat)
// jeszcze nie stoi. Dlatego strona ma dwa stany prawdziwe, nie jeden: „nie ma
// wydania" oraz „wydanie jest, kanał jeszcze nie stoi". Drugi stan bierze się
// z sekcji `kanal` w `wydania.json` i mówi się go wprost, przed wykazem plików.
import { tekst } from './szkielet.mjs';

// Reguła protokołu jest wymuszana, nie deklarowana. `wydania.json` mówi o polu
// `plik`: „adres https:// do pliku wydania — inny protokół jest odmawiany". Bez
// tego wymuszenia generator wystawiłby odnośnik o dowolnym protokole — także
// `javascript:` i `http:` — na jedynej publicznej stronie produktu.
//
// Odmowa opisuje brak. Wydanie z adresem innego protokołu nie znika ze strony
// i nie jest cenzurowane — traci odnośnik i dostaje zdanie, czego mu brakuje.
//
// REGUŁA PRZESZŁA PRÓBĘ 18.08.2026 I ZOSTAJE. Kanał tymczasowy na maszynie
// budującej stoi na samym adresie liczbowym, dla którego urząd certyfikacji
// certyfikatu nie wystawi — pierwszym odruchem było więc dopuszczenie tutaj
// `http:`. Wybrano drugie wyjście: postawiono TLS z certyfikatem własnym.
// Rozluźnienie tej stałej byłoby zmianą na zawsze, wprowadzoną dla wygody
// jednego dnia; certyfikat własny jest kosztem widocznym (ostrzeżenie
// przeglądarki) i odkręcalnym jednym przebiegiem `certbot --nginx`.
const PROTOKOL_WYDANIA = 'https:';

/** Adres, względem którego rozwijane są adresy względne w polu `plik`. Bierze się
 *  z `kanal.adres`, bo kanał wolno przenieść — a wpisana na stałe nazwa hosta
 *  rozwijałaby po przenosinach adresy względne względem hosta, którego już nie ma. */
function podstawaAdresu(wykaz) {
  const adres = wykaz && wykaz.kanal && typeof wykaz.kanal.adres === 'string' ? wykaz.kanal.adres : '';
  try {
    return new URL(adres).href;
  } catch {
    return 'https://danaco-group.pl/';
  }
}

/**
 * Adres do pobrania albo `null` wraz z powodem odmowy.
 *
 * Adres względny też przechodzi (wydanie leżące obok strony), bo rozstrzyga
 * protokół rozwiniętego adresu, a ten dziedziczy wtedy https po witrynie.
 */
function adresWydania(w, podstawa) {
  if (!w.plik) return { adres: null, powod: 'wydanie nie wskazuje pliku' };
  let rozwiniety;
  try {
    rozwiniety = new URL(w.plik, podstawa);
  } catch {
    return { adres: null, powod: 'adres pliku jest nieczytelny' };
  }
  if (rozwiniety.protocol !== PROTOKOL_WYDANIA) {
    return {
      adres: null,
      powod: `adres o protokole ${rozwiniety.protocol} — wykaz wydań przyjmuje wyłącznie https://`,
    };
  }
  return { adres: rozwiniety.href, powod: '' };
}

/**
 * Najnowsze wydanie wykazu — wyliczone, nie wzięte z pierwszej pozycji.
 *
 * Aplikacja (`client/src/aktualizacja/wykaz-wydan.ts`) kolejności nie ufa
 * i przechodzi cały wykaz porównując człony wersji liczbowo. Branie pozycji
 * zerowej dałoby przy wykazie z `1.0.0` przed `1.10.0` inny wynik niż aplikacja.
 * Ta sama reguła musi stać po obu stronach; przepisana jest tutaj, bo witryna nie
 * ma łańcucha budowy i nie zaciąga TypeScriptu klienta.
 */
function czlonyWersji(wersja) {
  return String(wersja ?? '')
    .trim()
    .replace(/^[vV]/, '')
    .split('.')
    .map((c) => (/^\d+$/.test(c.trim()) ? Number.parseInt(c.trim(), 10) : 0));
}

function nowszaWersja(kandydat, biezaca) {
  const a = czlonyWersji(kandydat);
  const b = czlonyWersji(biezaca);
  for (let i = 0; i < Math.max(a.length, b.length); i += 1) {
    const x = a[i] ?? 0;
    const y = b[i] ?? 0;
    if (x !== y) return x > y;
  }
  return false;
}

export function najnowszeWydanie(wydania) {
  let najnowsze = wydania[0];
  for (const pozycja of wydania) {
    if (nowszaWersja(pozycja.wersja, najnowsze.wersja)) najnowsze = pozycja;
  }
  return najnowsze;
}

/**
 * Wszystkie pozycje najnowszej wersji, nie jedna. Wydanie ma tyle plików, ile
 * systemów, a przy równych wersjach `najnowszeWydanie()` oddaje tylko pierwszą
 * pozycję tablicy — jeden plik z brzegu, niezależnie od systemu czytającego.
 * Ramka ma pokazać je wszystkie, więc bierze cały komplet najnowszej wersji.
 */
export function pozycjeNajnowszej(wydania) {
  const najnowsze = najnowszeWydanie(wydania);
  return wydania.filter((w) => !nowszaWersja(najnowsze.wersja, w.wersja));
}

/** Stan kanału pobierania, odczytany z wykazu. Brak sekcji `kanal` czytamy
 *  jako kanał niewdrożony: milczenie pliku nie jest deklaracją, że coś stoi. */
function stanKanalu(wykaz) {
  const kanal = wykaz && typeof wykaz.kanal === 'object' && wykaz.kanal !== null ? wykaz.kanal : {};
  return {
    wdrozony: kanal.wdrozony === true,
    // Kanał tymczasowy jest osobnym stanem, nie odcieniem wdrożonego. Pliki
    // odpowiadają — ale adres jest z maszyny budującej i nie ma być zapamiętany.
    tymczasowy: kanal.tymczasowy === true,
    czegoBrakuje: typeof kanal.czego_brakuje === 'string' ? kanal.czego_brakuje : '',
    podstawa: podstawaAdresu(wykaz),
  };
}

/**
 * Zdanie o niewdrożonym kanale — nad plikami, nie pod nimi. Pakiety istnieją, ale
 * adres, pod którym mają stanąć, może jeszcze nie odpowiadać. Operator musi to
 * przeczytać, zanim zobaczy nazwy plików, a nie po przewinięciu strony.
 */
function ostrzezenieKanalu(stan) {
  if (stan.wdrozony) return uwagaKanaluTymczasowego(stan);
  const brakuje = stan.czegoBrakuje ? ` Do uruchomienia kanału brakuje: ${tekst(stan.czegoBrakuje)}.` : '';
  return `  <div class="dc-uwaga dc-uwaga--brak">
    <p><strong>Wydanie jest zbudowane, ale kanał pobierania jeszcze nie stoi.</strong>
    Pliki poniżej istnieją naprawdę i ich sumy SHA-256 są policzone z tych właśnie plików —
    natomiast adres, pod którym mają stanąć, jeszcze nie odpowiada.${brakuje}</p>
    <p>Dlatego <strong>na tej stronie nie ma dziś przycisku „Pobierz"</strong>: przycisk, który nic
    nie odda, byłby kłamstwem. Zamiast niego stoi docelowy adres każdego pliku wraz z sumą —
    do sprawdzenia w dniu, w którym kanał ruszy. Do tego czasu pakiet wydaje się z ręki:
    <a href="mailto:support@danaco-group.pl">support@danaco-group.pl</a>.</p>
  </div>`;
}

/**
 * Zdanie o kanale, który STOI, ale jest tymczasowy.
 *
 * Kanał działający nie znaczy „nic więcej do powiedzenia": adres może być
 * tymczasowy, a to zaskoczyłoby pobierającego w połowie czynności, więc stoi
 * PRZED przyciskami, nie za nimi. Zdanie jest warunkowe i bierze się z wykazu:
 * gdy kanał dostał własną nazwę i certyfikat urzędu, znika samo, bo przestało
 * być prawdą. O hasło pyta nie kanał, tylko ŚCIEŻKA jednej pozycji — dlatego
 * zdanie o haśle stoi przy pozycji, nie tutaj.
 */
function uwagaKanaluTymczasowego(stan) {
  if (!stan.tymczasowy) return '';
  return `  <div class="dc-uwaga dc-uwaga--brak">
    <p><strong>Ten adres jest tymczasowy i przeglądarka pokaże ostrzeżenie o certyfikacie.</strong>
    Pliki stoją na maszynie budującej, pod adresem liczbowym, a dla adresu liczbowego urząd
    certyfikacji certyfikatu nie wystawia — połączenie jest więc szyfrowane certyfikatem
    podpisanym samym sobą i trzeba to ostrzeżenie świadomie przejść. Szyfrowanie jest tu
    warunkiem, nie ozdobą: bez niego hasło do pobierania szłoby przez sieć w postaci jawnej.
    Po przeniesieniu witryny na docelowy adres ostrzeżenie zniknie.</p>
    <p>Niezależnie od adresu i certyfikatu sprawdzianem, czy plik jest tym właściwym,
    zostaje <strong>suma SHA-256</strong> — stoi przy każdej pozycji niżej.</p>
  </div>`;
}

/** Nazwa systemu w postaci porównywalnej ze wskazaniem `navigator`. */
function kluczSystemu(w) {
  return String(w.system ?? '')
    .trim()
    .toLowerCase();
}

/**
 * Jeden pakiet w ramce wiodącej — z sumą przy pliku, nie tylko w tabeli niżej,
 * żeby Operator miał czym sprawdzić, co pobrał, bez szukania w chronologii.
 */
/** Rozmiar dla człowieka wraz z dokładną liczbą bajtów — ta druga jest tym, co
 *  pobierający porówna z nagłówkiem `content-length`, gdy zechce sprawdzić,
 *  czy plik doszedł cały. „48,6 MB" samo w sobie tego nie rozstrzyga. */
function rozmiarPozycji(w) {
  const ludzki = w.rozmiar ? String(w.rozmiar) : '';
  const bajty = Number.isFinite(w.rozmiarBajty)
    ? `${String(w.rozmiarBajty).replace(/\B(?=(\d{3})+(?!\d))/g, ' ')} B`
    : '';
  if (ludzki && bajty) return `${ludzki} (${bajty})`;
  return ludzki || bajty || 'rozmiar nieustalony';
}

/** Nagłówek karty. `nazwa` mówi system, architekturę i postać naraz; `system`
 *  sam z siebie nie odróżnia pozycji, bo wszystkie klienckie są „Windows".
 *  Brak `nazwy` cofa się do `system`, więc starszy wykaz nadal działa. */
function naglowekPozycji(w) {
  return w.nazwa ? String(w.nazwa) : String(w.system ?? 'System nieokreślony');
}

/**
 * Zdanie o haśle — przy pozycji, nie przy kanale. Uwierzytelnienie zamyka jedną
 * ścieżkę kanału, a nie kanał: pozycje spod `/pliki/` pobiera się bez pytania.
 * Zdanie postawione nad wszystkimi kartami stoi także nad tymi, które hasła nie
 * wymagają, i zapowiada okienko, które przy nich nie wyskakuje.
 */
function zdanieOHasle(w) {
  if (w.chronione_haslem !== true) return '';
  return `\n      <p class="dc-pakiet__haslo">Zapyta o hasło — ten plik leży pod ścieżką zamkniętą
      uwierzytelnieniem, a poświadczenie wydaje producent. Czytanie strony hasła nie wymaga.</p>`;
}

/**
 * Zdanie o braku podpisu — przy KAŻDEJ pozycji Windows, nie raz na całej stronie.
 * Pobierający czyta kartę tej pozycji, którą bierze, i przy niej ma się dowiedzieć,
 * że SmartScreen ukryje przycisk uruchomienia. Pole `podpisany` jest mierzone przy
 * składaniu instalki (katalog Security pliku), więc zdanie zniknie samo w dniu,
 * w którym plik dostanie podpis — nikt nie będzie go stąd wykreślał ręcznie.
 */
function zdanieOPodpisie(w) {
  if (String(w.system ?? '').toLowerCase() !== 'windows') return '';
  if (w.podpisany === true) return '';
  return `\n      <p class="dc-pakiet__brak">Plik nie jest podpisany — Windows SmartScreen pokaże
      ostrzeżenie o nieznanym wydawcy. Sprawdzianem, czy plik jest tym właściwym, zostaje suma niżej.</p>`;
}

function kartaPakietu(w, stan) {
  const { adres, powod } = adresWydania(w, stan.podstawa);
  const braki = [];
  if (!adres) braki.push(powod);
  if (!w.suma) braki.push('wydanie nie podaje sumy SHA-256, więc aplikacja go nie założy');

  let wiersz;
  if (braki.length > 0) {
    wiersz = `<p class="dc-pakiet__brak"><strong>Nie ma stąd czego bezpiecznie pobrać:</strong>
      ${tekst(braki.join('; '))}.</p>`;
  } else if (!stan.wdrozony) {
    // Świadomie bez znacznika <a>. Adres jest prawdziwy i docelowy, ale dziś
    // nie odpowiada — odnośnik udawałby, że coś odda. Adres stoi więc jako
    // tekst do skopiowania, opisany tym, czym jest.
    wiersz = `<p class="dc-pakiet__brak">Adres docelowy (jeszcze nie odpowiada):<br>
      <code>${tekst(adres)}</code></p>`;
  } else {
    wiersz = `<p><a class="btn btn--primary" href="${tekst(adres)}">Pobierz <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 3v12M6 13l6 6 6-6M4 21h16"/></svg></a></p>${zdanieOHasle(w)}`;
  }

  const coBierzesz = w.coBierzesz ? `\n      <p>${tekst(w.coBierzesz)}</p>` : '';

  return `    <div class="dc-pakiet" data-system="${tekst(kluczSystemu(w))}">
      <h3>${tekst(naglowekPozycji(w))}</h3>
      <p class="dc-pakiet__plik">${tekst(w.nazwaPliku ?? w.plik ?? '—')} — ${tekst(rozmiarPozycji(w))}</p>${coBierzesz}
      ${wiersz}${zdanieOPodpisie(w)}
      <p class="dc-pakiet__suma">SHA-256:<br><code>${tekst(w.suma ?? 'BRAK — aplikacja takiego wydania NIE ZAŁOŻY')}</code></p>
    </div>`;
}

/**
 * Rozpoznanie systemu Operatora — wyróżnienie, nie warunek widoczności. Nic, co
 * Operator ma zobaczyć, nie powstaje dopiero ze skryptu: wszystkie pakiety są
 * w HTML-u od pierwszej klatki, a skrypt jedynie dokłada znacznik przy tym, który
 * pasuje do czytającego, i zdanie o tym, co rozpoznał. Z wyłączonym JavaScriptem
 * strona jest kompletna — tylko bez wyróżnienia. Rozpoznanie idzie tym samym
 * wzorcem, co `systemBiezacy()` w kliencie
 * (`client/src/aktualizacja/wykaz-wydan.ts`), żeby strona i baner nie mówiły
 * dwóch różnych rzeczy o tej samej maszynie.
 */
const SKRYPT_ROZPOZNANIA = `
(function () {
  var opis = ((navigator.userAgent || '') + ' ' + (navigator.platform || '')).toLowerCase();
  var system = '';
  if (opis.indexOf('windows') >= 0 || opis.indexOf('win32') >= 0 || opis.indexOf('win64') >= 0) system = 'windows';
  if (!system) return;
  var pasujace = document.querySelectorAll('.dc-pakiet[data-system="' + system + '"]');
  if (pasujace.length === 0) return;
  var nazwa = 'Windows';
  for (var i = 0; i < pasujace.length; i += 1) {
    pasujace[i].classList.add('dc-pakiet--twoj');
    var znacznik = document.createElement('p');
    znacznik.className = 'dc-pakiet__znacznik';
    znacznik.textContent = 'Pasuje do twojego systemu (' + nazwa + ')';
    pasujace[i].insertBefore(znacznik, pasujace[i].firstChild);
  }
  var zdanie = document.getElementById('rozpoznanie');
  if (zdanie) {
    zdanie.textContent = 'Rozpoznany system: ' + nazwa + ' — pasujące pakiety są niżej wyróżnione. ' +
      'Pozostałe zostają widoczne, bo pobiera się też dla innej maszyny.';
    zdanie.hidden = false;
  }
})();
`;

/** Wiersz jednego wydania w chronologii. */
function wiersz(w, stan) {
  const { adres, powod } = adresWydania(w, stan.podstawa);
  let plik;
  if (!w.plik) {
    plik = '—';
  } else if (!adres) {
    plik = `${tekst(w.nazwaPliku ?? w.plik)} — <strong>bez odnośnika</strong>: ${tekst(powod)}`;
  } else if (!stan.wdrozony) {
    plik = `${tekst(w.nazwaPliku ?? w.plik)} — <strong>adres jeszcze nie odpowiada</strong>: <code>${tekst(adres)}</code>`;
  } else {
    plik = `<a href="${tekst(adres)}">${tekst(w.nazwaPliku ?? w.plik)}</a>`;
  }
  return `      <tr>
        <td><strong>${tekst(w.wersja)}</strong></td>
        <td>${tekst(w.data)}</td>
        <td>${tekst(naglowekPozycji(w))}</td>
        <td>${plik}</td>
        <td>${tekst(rozmiarPozycji(w))}</td>
        <td><code>${tekst(w.suma ?? 'BRAK — aplikacja takiego wydania NIE ZAŁOŻY')}</code></td>
        <td>${tekst(w.zmiany ?? '')}</td>
      </tr>`;
}

/**
 * Ramka wiodąca „najnowsze wydanie".
 *
 * Ramka wiodąca nie może być łagodniejsza niż chronologia. Odnośnik dostaje
 * wyłącznie wydanie z adresem https oraz sumą, i tylko przy stojącym kanale; brak
 * któregokolwiek z tych warunków jest tu nazwany, a nie przemilczany.
 */
function ramkaNajnowszego(wydania, stan) {
  const pozycje = pozycjeNajnowszej(wydania);
  const w = pozycje[0];
  return `  <h2>Najnowsze wydanie: ${tekst(w.wersja)} z dnia ${tekst(w.data)}</h2>
  <p id="rozpoznanie" class="dc-rozpoznanie" hidden></p>
  <div class="dc-pakiety">
${pozycje.map((p) => kartaPakietu(p, stan)).join('\n')}
  </div>
  <script>${SKRYPT_ROZPOZNANIA}</script>`;
}

/**
 * Pozycje zapowiedziane, których pliku jeszcze nie ma.
 *
 * PO CO OSOBNY KLUCZ, A NIE PUSTA POZYCJA W WYKAZIE. Pozycja bez pliku i bez sumy
 * w kluczu `wydania` byłaby widziana przez baner aktualizacji w aplikacji, który
 * czyta ten sam plik — a baner nie ma prawa proponować czegoś, czego nie ma.
 * Milczenie o pakiecie, który się buduje, byłoby jednak drugą nieprawdą: ktoś
 * szukający zapowiedzianego pakietu zobaczyłby sam wykaz i wyszedł z przekonaniem,
 * że nic więcej nie powstaje. Dlatego jest tu miejsce nazwane wprost
 * „w przygotowaniu" — bez rozmiaru i BEZ SUMY, bo suma zmyślona jest najgorszą
 * możliwą nieprawdą na tej stronie: unieważnia jedyny sprawdzian, jaki tu stoi.
 */
function wPrzygotowaniu(wykaz) {
  const pozycje = Array.isArray(wykaz.w_przygotowaniu) ? wykaz.w_przygotowaniu : [];
  if (pozycje.length === 0) return '';
  const karty = pozycje
    .map(
      (w) => `    <div class="dc-pakiet dc-pakiet--przygotowanie">
      <h3>${tekst(naglowekPozycji(w))}</h3>
      <p class="dc-pakiet__plik">${tekst(w.nazwaPliku ?? '—')}</p>
      ${w.coBierzesz ? `<p>${tekst(w.coBierzesz)}</p>` : ''}
      <p class="dc-pakiet__brak"><strong>W przygotowaniu — pliku jeszcze nie ma.</strong>
      Nie ma go pod żadnym adresem, więc nie ma tu przycisku ani sumy kontrolnej.
      Pozycja pojawi się w wykazie wyżej w tej samej chwili, w której plik powstanie
      i zostanie zmierzony.</p>
    </div>`,
    )
    .join('\n');
  return `  <h2>W przygotowaniu</h2>
  <div class="dc-pakiety">
${karty}
  </div>`;
}

/**
 * Pakiet serwera — ODDZIELONY OD POZYCJI KLIENCKICH, i to celowo.
 *
 * Leży w osobnym kluczu wykazu, nie w tablicy `wydania`: nie jest wydaniem dla
 * klienta i nie ma się nigdy pokazać w banerze aktualizacji aplikacji. Na stronie
 * stoi pod własnym nagłówkiem, poniżej pozycji klienckich, z jednym zdaniem
 * o tym, dla kogo jest. Klient, który weźmie ten plik zamiast swojego, straci
 * wieczór — więc różnica musi być widoczna, nie domyślna.
 */
function sekcjaSerwera(wykaz, stan) {
  const s = wykaz && typeof wykaz.serwer === 'object' && wykaz.serwer !== null ? wykaz.serwer : null;
  if (!s) return '';
  return `  <h2>Dla administratora: pakiet serwera wdrożenia</h2>
  <div class="dc-uwaga">
    <p><strong>To nie jest pozycja dla klienta.</strong> ${tekst(s.czym_to_jest ?? '')}</p>
  </div>
  <div class="dc-pakiety">
${
  // `system: 'serwer'` podmienia się świadomie: pakiet administratora nie ma
  // nigdy trafić pod wzorzec skryptu rozpoznającego maszynę czytającego, choćby
  // wzorzec kiedyś rozszerzono. Nagłówek karty bierze się z `nazwa`, więc
  // informacja o systemie i architekturze nie ginie.
  kartaPakietu({ ...s, system: 'serwer' }, stan)
}
  </div>`;
}

/** Buduje treść strony „Pobierz" z wykazu wydań. */
export function stronaPobierz(strona, wykaz) {
  const wydania = Array.isArray(wykaz.wydania) ? wykaz.wydania : [];
  const stan = stanKanalu(wykaz);

  // Nazwa pliku w przykładzie polecenia bierze się z wykazu, nie z pamięci:
  // wpisana na stałe zestarzała się przy pierwszym wydaniu i uczyła pobierającego
  // przepisywać nazwę cudzego pliku. Bez wydania przykładu nie ma wcale —
  // polecenie z wymyśloną nazwą nie zadziała u nikogo.
  const przykladSumy =
    wydania.length === 0
      ? ''
      : `<br>\n  <code>certutil -hashfile &quot;${tekst(pozycjeNajnowszej(wydania)[0].nazwaPliku ?? '')}&quot; SHA256</code>`;

  const najnowsze =
    wydania.length === 0
      ? `  <div class="dc-uwaga">
    <p><strong>Nie ma jeszcze żadnego wydania.</strong> Nie istnieje plik, który dałoby się
    stąd pobrać. Ta strona pokaże wydanie w tej samej chwili, w której ono powstanie —
    a do tego czasu mówi o tym wprost, zamiast pokazywać przycisk bez pliku.</p>
  </div>`
      : ramkaNajnowszego(wydania, stan);

  const chronologia =
    wydania.length === 0
      ? `<p>Wykaz jest pusty. Każde wydanie — także poprawkowe — zostanie tu dopisane
  wraz z datą i pozostanie na stałe: chronologia ma pokazywać, co i kiedy się zmieniało,
  a nie wyłącznie stan bieżący.</p>`
      : `<div class="dc-tabela"><table>
    <thead>
      <tr><th>Wersja</th><th>Data</th><th>System</th><th>Plik</th><th>Rozmiar</th><th>Suma SHA-256</th><th>Co się zmieniło</th></tr>
    </thead>
    <tbody>
${wydania.map((w) => wiersz(w, stan)).join('\n')}
    </tbody>
  </table></div>`;

  return {
    ...strona,
    tresc: `  <h1>Pobierz Danaco Console</h1>
  <p class="wiodacy">Instalka aplikacji oraz pełna chronologia aktualizacji — każde wydanie z datą.</p>

${wydania.length === 0 ? '' : ostrzezenieKanalu(stan)}

${najnowsze}

${wPrzygotowaniu(wykaz)}

${sekcjaSerwera(wykaz, stan)}

  <h2>Co dostajesz w instalce</h2>
  <p>Produkt występuje w <strong>jednej postaci — hybrydzie</strong>. Instalka niesie samo okno
  wraz z powłoką; rdzeń — czyli to, co naprawdę pracuje — <strong>nie jest w niej w ogóle</strong>.
  Rdzeń stoi na serwerze wdrożenia, a okno się z nim łączy. Pakiet jest przez to mały,
  a przetwarzanie multimediów i dokumentów odbywa się na serwerze, nie na stanowisku.</p>
  <p><strong>Bez działającego serwera wdrożenia okno nie ma z czym rozmawiać</strong> i pokaże to
  od pierwszego uruchomienia. Serwer zakłada administrator raz, z pakietu stojącego wyżej na tej
  stronie; stanowiska dostają samo okno.</p>

  <h2>Numer wersji i data</h2>
  <p>Wydania różni <strong>data</strong>, nie sam numer: pliki z różnych dni leżą w osobnych
  katalogach z datą w nazwie i nie nadpisują się wzajemnie. Numer podnosi się rozstrzygnięciem
  producenta; bieżący stoi wyżej, przy najnowszym wydaniu, i tam też jest data.</p>

  <h2>Windows: x64 czy ARM64</h2>
  <p>Zwykły komputer i laptop z procesorem Intel albo AMD to <strong>x64</strong> — ta pozycja
  pasuje niemal każdemu. <strong>ARM64</strong> jest dla maszyn na procesorze ARM (np. Surface
  na Snapdragonie). Plik ARM64 nie uruchomi się na maszynie Intel/AMD ani odwrotnie. Sprawdzenie:
  Ustawienia Windows → System → Informacje → „Typ systemu".</p>

  <h2>Podpis kodu</h2>
  <p>Przy każdej pozycji Windows stoi zdanie o jej podpisie. Plik bez podpisu daje przy pierwszym
  uruchomieniu ostrzeżenie SmartScreen o nieznanym wydawcy — nie jest to usterka instalki i nie
  znaczy, że plik jest podmieniony; znaczy, że wydawca nie ma jeszcze certyfikatu podpisującego.
  Sposób na sprawdzenie, czy plik jest tym właściwym, jest wtedy jeden i stoi przy pozycji:
  suma SHA-256.</p>

  <h2>Aktualizacja z poziomu aplikacji</h2>
  <p>Zainstalowana aplikacja sama sprawdza, czy stoi tu coś nowszego. Gdy tak jest,
  w oknie pojawia się baner <strong>„Aktualizuje"</strong> działający jak przycisk:
  jedno kliknięcie pobiera wydanie, zakłada je i uruchamia aplikację ponownie.
  Dopóki kanał pobierania nie stoi, aplikacja nie ma czego odczytać i baner się
  nie pokazuje — nie jest to usterka, tylko brak kanału.</p>

  <h2>Suma kontrolna jest warunkiem, nie ozdobą</h2>
  <p>Każde wydanie niesie sumę SHA-256. Aplikacja aktualizująca się banerem liczy ją
  z pobieranego strumienia i <strong>odmawia założenia pliku, którego suma się nie zgadza</strong> —
  plik jest wtedy kasowany, zanim cokolwiek zostanie podmienione. Wydanie bez podanej sumy
  nie zostanie założone w ogóle. Pobierając ręcznie, porównaj sumę samodzielnie — na Windows
  robi to narzędzie systemowe, bez zakładania czegokolwiek:${przykladSumy}</p>
  <p>Certyfikat kanału potwierdza, z kim rozmawia przeglądarka. Suma kontrolna potwierdza coś
  innego i dlatego jedno nie zastępuje drugiego: że <strong>ten konkretny plik</strong> jest tym,
  który producent zbudował — niezależnie od tego, jaką drogą przyszedł i ile razy był kopiowany.</p>`,
  };
}
