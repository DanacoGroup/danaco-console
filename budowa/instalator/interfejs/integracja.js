/* Skrypt spoza warstwy projektowej: wiąże kroki 1, 4 i 5 kreatora z rzeczywistym
   stanem tej maszyny. Nie zmienia żadnego pliku design/ — dokłada własne węzły
   do już zmontowanego okna (zdarzenie `kreator-gotowy`) i czyta dane, które
   main.rs złożył w adresie okna albo udostępnił poleceniami Tauri. Wartość,
   której nie da się odczytać, jest tu nazwaną odmową w tekście węzła — nigdy
   liczbą ani ścieżką zmyśloną. */
(function () {
'use strict';

var TAURI = window.__TAURI__;
var invoke = TAURI && TAURI.core && TAURI.core.invoke;
var ADRES = new URLSearchParams(location.search);

/* Podmiana katalogu treści musi wykonać się PRZED zdarzeniem `kreator-gotowy`:
   `okna/instalator.js` czyta krok5.blad.podtytul w chwili tego zdarzenia i
   zapamiętuje go na cały czas życia okna. Domyślny podtytuł zakłada niepowodzenie
   akurat na etapie rejestrowania składników — zdanie nieprawdziwe dla odmowy
   ustalonej tu, w main.rs, przed otwarciem okna (krok 5 zatrzymuje się przed
   pobraniem, nie w jego trakcie). Reszta bloku błędu (kod i rada) zostaje
   nadpisana później, w DOM, po zdarzeniu — tamte węzły `ustawOdslone` nie dotyka. */
(function podmienPodtytulBledu() {
  var K = window.DanacoKreator;
  if (!K || !K.tresci || !K.tresci.krok5 || !K.tresci.krok5.blad) return;
  K.tresci.krok5.blad.podtytul =
    'Instalacja zatrzymała się, zanim ukończyła pobranie i zapis składników programu. '
    + 'Zmiany zostały cofnięte — w komputerze nie pozostały pliki programu.';
})();

function el(znacznik, klasa, tekst) {
  var n = document.createElement(znacznik);
  if (klasa) n.className = klasa;
  if (tekst !== undefined) n.textContent = tekst;
  return n;
}

function wiersz(etykieta, wartosc) {
  var w = el('div', 'dn-zestawienie dn-zestawienie--cichy');
  w.appendChild(el('span', null, etykieta));
  w.appendChild(el('b', null, wartosc));
  return w;
}

/* Ten sam układ „nagłówek + dn-wykaz-cichy" co blokDanych z biblioteki
   składników, złożony wprost — blokDanych i naglowekBloku czytają katalog
   treści po ścieżce, a treść tu jest odczytem maszyny, nie kluczem katalogu. */
function blokOdczytu(tytul, wiersze) {
  var oprawa = el('div');
  oprawa.appendChild(el('h3', 'dn-kreator-podtytul', tytul));
  var wykaz = el('div', 'dn-wykaz-cichy');
  wiersze.forEach(function (w) { wykaz.appendChild(w); });
  oprawa.appendChild(wykaz);
  return oprawa;
}

function formatujBajty(bajty) {
  if (bajty == null) return null;
  var gb = bajty / 1e9;
  if (gb >= 1) return gb.toFixed(1).replace('.', ',') + ' GB';
  return Math.round(bajty / 1e6) + ' MB';
}

function mapaProcesora(kod) {
  return { x64: 'Intel/AMD (x64)', arm: 'ARM64', brak: 'nierozpoznana' }[kod] || kod;
}

function odmowa(tekst) { return 'Odmowa: ' + (tekst || 'nie udało się odczytać'); }

/* ——— Krok 1: stan tej maszyny, doklejony pod wymaganiami produktu ————— */
function wiazKrok1(stan) {
  var sekcja = document.querySelector('.dn-kreator-ekran[data-ekran="1"]');
  var istniejacy = sekcja && sekcja.querySelector('.dn-wykaz-cichy');
  if (!istniejacy || !istniejacy.parentElement) return;
  var blok = blokOdczytu('Stan tego urządzenia (odczyt maszyny)', [
    wiersz('Wersja systemu', stan.system),
    wiersz('Architektura procesora', mapaProcesora(stan.procesor)),
    wiersz('Wolne miejsce na dysku', stan.wolne_miejsce_bajty != null
      ? formatujBajty(stan.wolne_miejsce_bajty) : odmowa(stan.wolne_miejsce_odmowa)),
    wiersz('Uprawnienia procesu instalatora', stan.uprawnienia)
  ]);
  istniejacy.parentElement.insertAdjacentElement('afterend', blok);
}

/* ——— Krok 4: ścieżki domyślne tej maszyny i sprawdzenie zapisu próbą ———— */
function wiazKrok4(stan) {
  var sekcja = document.querySelector('.dn-kreator-ekran[data-ekran="4"]');
  if (!sekcja) return;

  ustawPole('katalog-programu', stan.katalog_programu, stan.katalog_programu_odmowa);
  ustawPole('katalog-danych', stan.katalog_danych, stan.katalog_danych_odmowa);

  var miejsce = sekcja.querySelector('.dn-meta');
  var blokSprawdzenia = blokOdczytu('Sprawdzenie zapisu (odczyt maszyny)', [
    wiersz('Prawo zapisu', 'sprawdzanie…'),
    wiersz('Wolne miejsce na tym dysku', 'sprawdzanie…')
  ]);
  if (miejsce && miejsce.parentElement) miejsce.insertAdjacentElement('afterend', blokSprawdzenia);
  else sekcja.appendChild(blokSprawdzenia);

  if (stan.katalog_programu) {
    odswiezSprawdzenie(stan.katalog_programu, blokSprawdzenia);
  } else {
    ustawWynikSprawdzenia(blokSprawdzenia, {
      zapisywalny: false,
      powod_odmowy: stan.katalog_programu_odmowa,
      wolne_miejsce_bajty: null
    });
  }

  wiazZmienPrzycisk('katalog-programu', blokSprawdzenia);
  wiazZmienPrzycisk('katalog-danych', null);
}

function ustawPole(id, wartosc, odmowaTekst) {
  var pole = document.getElementById(id);
  if (pole) pole.value = wartosc || odmowa(odmowaTekst);
}

function przyciskPola(id) {
  var pole = document.getElementById(id);
  var zestaw = pole && pole.closest('.dn-pole-zestaw');
  return zestaw ? zestaw.querySelector('button') : null;
}

function wiazZmienPrzycisk(id, blokSprawdzenia) {
  var przycisk = przyciskPola(id);
  if (!przycisk || !invoke) return;
  przycisk.addEventListener('click', function () {
    invoke('wybierz_katalog', { tytul: 'Wskaż katalog' }).then(function (wynik) {
      if (!wynik) return;
      var pole = document.getElementById(id);
      if (pole) pole.value = wynik.sciezka;
      if (blokSprawdzenia) ustawWynikSprawdzenia(blokSprawdzenia, wynik);
    });
  });
}

function odswiezSprawdzenie(sciezka, blok) {
  if (!invoke) {
    ustawWynikSprawdzenia(blok, { zapisywalny: false, powod_odmowy: 'środowisko Tauri niedostępne', wolne_miejsce_bajty: null });
    return;
  }
  invoke('sprawdz_katalog', { sciezka: sciezka }).then(function (wynik) {
    ustawWynikSprawdzenia(blok, wynik);
  });
}

function ustawWynikSprawdzenia(blok, wynik) {
  var wykaz = blok.querySelector('.dn-wykaz-cichy');
  wykaz.innerHTML = '';
  wykaz.appendChild(wiersz('Prawo zapisu', wynik.zapisywalny ? 'Tak — potwierdzone próbą zapisu' : odmowa(wynik.powod_odmowy)));
  wykaz.appendChild(wiersz('Wolne miejsce na tym dysku', wynik.wolne_miejsce_bajty != null
    ? formatujBajty(wynik.wolne_miejsce_bajty) : odmowa('nie udało się odczytać wolnego miejsca')));
}

/* ——— Krok 5: powód rzeczywistej odmowy pobrania, ustalony przed otwarciem
   okna (main.rs wykonuje jedno prawdziwe żądanie HTTPS do serwera wydań
   Danaco zanim okno stanie) — tu tylko wpisany w blok błędu już zmontowany
   przez ekrany/5-instalacja.js. Nie dotyka paska postępu: krok wchodzi w
   odsłonę błędu wprost, bez udawanego przebiegu, którego adres nie potwierdza. */
function wiazKrok5Blad() {
  var powod = ADRES.get('blad_powod');
  var zdanie = ADRES.get('blad_zdanie');
  if (!powod && !zdanie) return;
  var szczegoly = document.querySelector('[data-blad-szczegoly]');
  var rada = document.querySelector('.dn-kreator-ekran[data-ekran="5"] [data-blok="blad"] p');
  if (szczegoly) szczegoly.textContent = 'Powód: ' + powod + '\n' + zdanie;
  if (rada) rada.textContent = 'Uzyskaj poświadczenia dostępu do kanału wydań albo skontaktuj się '
    + 'z administratorem serwera wdrożenia, a następnie uruchom instalator ponownie.';
}


/* ——— Krok 5: rzeczywiste pobranie składników —————
   Kreator sam nie pobiera niczego — prototyp odgrywa przebieg czasem. Produkt
   przejmuje go tutaj: miara idzie z bajtów odebranych przez `pobierz_skladniki`,
   a krok 6 wchodzi dopiero po sprawdzeniu sumy pliku po stronie Rust. */
var ZDARZENIE_POSTEPU = 'instalator:postep-pobrania';

/** Katalog docelowy wskazany w kroku 4; pusty znaczy, że kroku nie wypełniono. */
function katalogDocelowy() {
  var pole = document.getElementById('katalog-programu');
  var wartosc = pole ? pole.value.trim() : '';
  return wartosc.indexOf('Odmowa:') === 0 ? '' : wartosc;
}

/* Rady dobrane do powodu odmowy. Rada prototypu mówi o zamknięciu innych
   programów i prawach administratora — dla braku mostu ani dla braku katalogu
   nie jest prawdą, a rada nieprawdziwa jest gorsza niż jej brak. */
var RADY = {
  'brak-mostu-programu': 'Uruchom pobrany plik instalacyjny Danaco Console.',
  'brak-katalogu': 'Wróć do kroku czwartego i wskaż katalog, w którym program ma stanąć.'
};

/** Wpisuje powód odmowy wraz z radą w blok błędu kroku 5 i oddaje go odsłonie błędu. */
function odmowaKroku5(bieg, powod, zdanie) {
  var szczegoly = document.querySelector('[data-blad-szczegoly]');
  if (szczegoly) szczegoly.textContent = 'Powód: ' + powod + '\n' + zdanie;
  var rada = document.querySelector('.dn-kreator-ekran[data-ekran="5"] [data-blok="blad"] p');
  if (rada && RADY[powod]) rada.textContent = RADY[powod];
  bieg.odmowa();
}

function przebiegKroku5(bieg) {
  if (!invoke) {
    odmowaKroku5(bieg, 'brak-mostu-programu',
      'Okno stoi poza powłoką programu, więc pobrania nie ma czym wykonać. '
      + 'Instalację prowadzi się wyłącznie z pliku instalacyjnego.');
    return;
  }
  var katalog = katalogDocelowy();
  if (!katalog) {
    odmowaKroku5(bieg, 'brak-katalogu',
      'Katalog docelowy nie został ustalony w kroku czwartym, więc nie ma dokąd zapisać plików.');
    return;
  }
  var odsluch = TAURI.event && TAURI.event.listen
    ? TAURI.event.listen(ZDARZENIE_POSTEPU, function (zdarzenie) {
        var d = zdarzenie.payload;
        if (d && d.razem_bajtow > 0) bieg.postep(d.odebrano_bajtow / d.razem_bajtow);
      })
    : null;
  function odlacz() { if (odsluch) odsluch.then(function (stop) { stop(); }); }

  invoke('stan_maszyny').then(function (stan) {
    return invoke('pobierz_skladniki', {
      architektura: stan.procesor,
      katalogRoboczy: katalog
    });
  }).then(function () {
    odlacz();
    bieg.postep(1);
    bieg.koniec();
  }).catch(function (blad) {
    odlacz();
    odmowaKroku5(bieg, (blad && blad.powod) || 'przebieg-przerwany',
      (blad && blad.zdanie) || String(blad));
  });
}

// Przejęcie kroku 5 wchodzi przed zdarzeniem gotowości: kreator czyta wpis
// w chwili wejścia w krok, a wejść może zaraz po zmontowaniu okna.
if (window.DanacoKreator) window.DanacoKreator.przebiegKroku5 = przebiegKroku5;

document.addEventListener('kreator-gotowy', function () {
  wiazKrok5Blad();
  if (!invoke) return;
  invoke('stan_maszyny').then(function (stan) {
    wiazKrok1(stan);
    wiazKrok4(stan);
  });
});
})();
