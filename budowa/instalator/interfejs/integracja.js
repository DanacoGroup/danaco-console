/* Skrypt spoza warstwy projektowej: wiąże kroki 1, 4, 5 i 6 kreatora z rzeczywistym
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
var K = window.DanacoKreator;

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

/* ——— Krok 4: ścieżki tej maszyny i sprawdzenie zapisu próbą ———————————— */
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
  /* Katalog danych składa sobie sam program z profilu Operatora i nie czyta
     znikąd innego wskazania, więc instalacja nie ma czym wyboru spełnić —
     pole zostaje odczytem, a czynność „Zmień…" schodzi z okna. */
  zdejmijZmienPrzycisk('katalog-danych');
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

function zdejmijZmienPrzycisk(id) {
  var przycisk = przyciskPola(id);
  if (przycisk) przycisk.remove();
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

/* ——— Krok 5: pobranie, sprawdzenie sumy i założenie programu ——————————
   Kreator sam nie instalował niczego — prototyp odgrywa przebieg czasem.
   Produkt przejmuje go tutaj: miara idzie z bajtów odebranych przez
   `pobierz_skladniki`, a krok 6 wchodzi dopiero wtedy, gdy instalka wydania
   uruchomiona przez `zaloz_program` odda kod zerowy. */
var ZDARZENIE_POSTEPU = 'instalator:postep-pobrania';
var ZDARZENIE_KANALU = 'instalator:stan-kanalu';

/* Miara kroku 5 jest miarą CAŁOŚCI, a okno rozdziela ją na cztery etapy
   udziałami 0,08 / 0,46 / 0,34 / 0,12 (`okna/instalator.js`). Progi poniżej
   trzymają etap wskazany na liście zgodnym z czynnością naprawdę wykonywaną. */
var MIARA_SPRAWDZENIE = 0.04;
var MIARA_POBIERANIE_OD = 0.08;
var MIARA_POBIERANIE_DO = 0.54;
var MIARA_ZAKLADANIE = 0.90;

/* Wynik sprawdzenia kanału wydań z main.rs. `null` znaczy „jeszcze nie wrócił",
   nie „kanał odpowiada": krok 5 rusza wtedy pobraniem, a odmowa serwera wychodzi
   z niego samego. */
var stanKanalu = null;

/* Założenie programu z kroku 5 — jedyne źródło ścieżki, którą uruchamia krok 6. */
var zalozone = null;

if (TAURI && TAURI.event && TAURI.event.listen) {
  TAURI.event.listen(ZDARZENIE_KANALU, function (zdarzenie) {
    stanKanalu = zdarzenie.payload || null;
  });
}

/** Katalog docelowy wskazany w kroku 4; pusty znaczy, że kroku nie wypełniono. */
function katalogDocelowy() {
  var pole = document.getElementById('katalog-programu');
  var wartosc = pole ? pole.value.trim() : '';
  return wartosc.indexOf('Odmowa:') === 0 ? '' : wartosc;
}

/** Wersja wybrana w kroku 3; wykrycie maszyny jest tam tylko zaznaczeniem domyślnym. */
function wybranaArchitektura() {
  var wybor = document.querySelector('[name="wariant"]:checked');
  if (wybor) return wybor.value;
  return ADRES.get('procesor') || 'brak';
}

/* Rady dobrane do powodu odmowy. Rada z katalogu treści mówi o ponownym
   uruchomieniu i o uprawnieniach administratora — dla braku mostu, braku
   katalogu i milczącego kanału nie jest prawdą, a rada nieprawdziwa jest
   gorsza niż jej brak. */
var RADY = {
  'brak-mostu-programu': 'Uruchom pobrany plik instalacyjny Danaco Console.',
  'brak-katalogu': 'Wróć do kroku czwartego i wskaż katalog, w którym program ma stanąć.',
  'brak-lacznosci': 'Sprawdź połączenie z internetem i uruchom instalację ponownie.',
  'odpowiedz-serwera': 'Uzyskaj poświadczenia dostępu do kanału wydań albo skontaktuj się '
    + 'z administratorem serwera wdrożenia, a następnie uruchom instalację ponownie.',
  'polaczenie-przerwane': 'Sprawdź połączenie z internetem i uruchom instalację ponownie.',
  'suma-niezgodna': 'Uruchom instalację ponownie. Jeżeli suma znów się nie zgodzi, pobierz '
    + 'instalator z kanału wydań na nowo.',
  'wydanie-nieopublikowane': 'Wróć do kroku trzeciego i wybierz wersję opublikowaną w kanale wydań.',
  'architektura-nierozpoznana': 'Wróć do kroku trzeciego i wskaż wersję programu.',
  'instalka-nie-ruszyla': 'Uruchom instalację ponownie jako administrator.',
  'instalka-odmowila': 'Zamknij otwarte okna Danaco Console i uruchom instalację ponownie.',
  'program-nieodnaleziony': 'Wróć do kroku czwartego i wskaż katalog, do którego wolno zapisywać.'
};

/** Wpisuje powód odmowy wraz z radą w blok błędu kroku 5 i oddaje go odsłonie błędu. */
function odmowaKroku5(bieg, powod, zdanie) {
  var szczegoly = document.querySelector('[data-blad-szczegoly]');
  if (szczegoly) szczegoly.textContent = 'Powód: ' + powod + '\n' + zdanie;
  var rada = document.querySelector('.dn-kreator-ekran[data-ekran="5"] [data-blok="blad"] p');
  if (rada && RADY[powod]) rada.textContent = RADY[powod];
  bieg.odmowa();
}

/** Cel licznika etapu pobierania — miara bierze się z rozmiaru pozycji wydania. */
function ustawCelPobrania(razemBajtow) {
  var etap = document.querySelector('[data-etapy] > .dn-krok[data-licznik-czasownik]');
  if (etap && razemBajtow > 0) etap.dataset.licznikCel = Math.round(razemBajtow / 1e6);
}

function przebiegKroku5(bieg) {
  if (!invoke) {
    odmowaKroku5(bieg, 'brak-mostu-programu',
      'Okno stoi poza powłoką programu, więc instalacji nie ma czym wykonać. '
      + 'Instalację prowadzi się wyłącznie z pliku instalacyjnego.');
    return;
  }
  var katalog = katalogDocelowy();
  if (!katalog) {
    odmowaKroku5(bieg, 'brak-katalogu',
      'Katalog docelowy nie został ustalony w kroku czwartym, więc nie ma dokąd zapisać plików.');
    return;
  }
  if (stanKanalu && stanKanalu.odmowa) {
    odmowaKroku5(bieg, stanKanalu.odmowa.powod, stanKanalu.odmowa.zdanie);
    return;
  }

  var odsluch = TAURI.event && TAURI.event.listen
    ? TAURI.event.listen(ZDARZENIE_POSTEPU, function (zdarzenie) {
        var d = zdarzenie.payload;
        if (!d || !(d.razem_bajtow > 0)) return;
        ustawCelPobrania(d.razem_bajtow);
        var udzial = Math.min(d.odebrano_bajtow / d.razem_bajtow, 1);
        bieg.postep(MIARA_POBIERANIE_OD + udzial * (MIARA_POBIERANIE_DO - MIARA_POBIERANIE_OD));
      })
    : null;
  function odlacz() { if (odsluch) odsluch.then(function (stop) { stop(); }); }

  bieg.postep(MIARA_SPRAWDZENIE);
  invoke('pobierz_skladniki', {
    architektura: wybranaArchitektura(),
    katalogRoboczy: katalog
  }).then(function (pobrane) {
    /* Sumę sprawdza Rust w obrębie pobrania: `pobierz_skladniki` kończy się
       dopiero po jej zgodności, więc etap sumy wchodzi tu już jako zamknięty. */
    opiszWiersz('wersja', pobrane.wersja);
    bieg.postep(MIARA_ZAKLADANIE);
    return invoke('zaloz_program', {
      plikInstalki: pobrane.plik_instalki,
      katalogProgramu: katalog
    });
  }).then(function (wynik) {
    odlacz();
    zalozone = wynik;
    opiszWiersz('lokalizacja', wynik.katalog_programu);
    bieg.postep(1);
    bieg.koniec();
  }).catch(function (blad) {
    odlacz();
    odmowaKroku5(bieg, (blad && blad.powod) || 'przebieg-przerwany',
      (blad && blad.zdanie) || String(blad));
  });
}

/* ——— Krok 6: podsumowanie mówi to, co naprawdę stanęło ————————————————
   Wiersze podsumowania rozpoznaje się po etykiecie z katalogu treści, bo blok
   danych nie znakuje ich atrybutem, a kolejność w bloku nie jest umową. */
function opiszWiersz(nazwaPola, wartosc) {
  var szczegoly = (K && K.tresci && K.tresci.krok6 && K.tresci.krok6.szczegoly) || {};
  var pole = szczegoly[nazwaPola];
  var ekran = document.querySelector('.dn-kreator-ekran[data-ekran="6"]');
  if (!pole || !pole.etykieta || !ekran || !wartosc) return;
  var wiersze = ekran.querySelectorAll('.dn-zestawienie');
  for (var i = 0; i < wiersze.length; i++) {
    var podpis = wiersze[i].querySelector('span');
    var miejsce = wiersze[i].querySelector('b');
    if (podpis && miejsce && podpis.textContent === pole.etykieta) {
      miejsce.textContent = wartosc;
      return;
    }
  }
}

/** Ogłoszenie kroku 6 wchodzi we frazę nawigacyjną — jedyny obszar `aria-live` tego ekranu. */
function oglosNaKroku6(zdanie) {
  var fraza = document.querySelector('.dn-kreator-ekran[data-ekran="6"] .dn-kreator-nawigacja');
  if (fraza) fraza.textContent = zdanie;
}

/* Czynność domyślna kroku 6 uruchamia program założony w kroku 5 i dopiero po
   jego starcie zamyka kreator. Ścieżka pochodzi z założenia — kreator nie składa
   jej drugi raz z katalogu i nie zgaduje nazwy pliku wykonywalnego. */
function uruchomProgram(poUruchomieniu) {
  if (!invoke || !zalozone) {
    oglosNaKroku6('Program nie został założony w tym przebiegu instalatora, '
      + 'więc nie ma czego uruchomić.');
    return;
  }
  invoke('uruchom_program', { sciezkaProgramu: zalozone.sciezka_programu })
    .then(function () { poUruchomieniu(); })
    .catch(function (blad) {
      oglosNaKroku6((blad && blad.zdanie) || String(blad));
    });
}

/* ——— Belka okna ————————————————————————————————————————————————————————
   Okno platformy stoi bez ramy, więc zwinięcie, rozwinięcie, zamknięcie
   i przeciąganie okna wykonuje belka kreatora. Przyciski rozpoznaje się po
   etykiecie z katalogu treści, bo kolejność w belce nie jest umową. */
function oknoPowloki() {
  return TAURI && TAURI.window && TAURI.window.getCurrentWindow
    ? TAURI.window.getCurrentWindow()
    : null;
}

function wiazBelke() {
  var okno = oknoPowloki();
  if (!okno) return;
  var belka = document.querySelector('.dn-kreator-belka');
  if (belka) belka.setAttribute('data-tauri-drag-region', '');
  var tresci = (K && K.tresci && K.tresci.okno) || {};
  var czynnosci = [[tresci.zwin, 'minimize'], [tresci.rozwin, 'toggleMaximize']];
  czynnosci.forEach(function (para) {
    if (!para[0]) return;
    var btn = document.querySelector('.dn-kreator-belka-btn[aria-label="' + para[0] + '"]');
    if (btn) btn.addEventListener('click', function () { okno[para[1]](); });
  });
}

/* Wyjście z okna: kreator pyta o zgodę, gdy jest o co, i dopiero wtedy woła
   ten zamek. Zamknięcie okna kończy proces instalatora. */
if (K) {
  K.wyjscieOkna = function () {
    var okno = oknoPowloki();
    if (okno) okno.close();
  };
}

// Przejęcie kroku 5 i czynności kroku 6 wchodzi przed zdarzeniem gotowości:
// kreator czyta oba wpisy w chwili wejścia w krok, a wejść może zaraz po
// zmontowaniu okna.
if (K) {
  K.przebiegKroku5 = przebiegKroku5;
  K.uruchomProgram = uruchomProgram;
}

document.addEventListener('kreator-gotowy', function () {
  wiazBelke();
  if (!invoke) return;
  invoke('stan_maszyny').then(function (stan) {
    wiazKrok1(stan);
    wiazKrok4(stan);
  });
});
})();
