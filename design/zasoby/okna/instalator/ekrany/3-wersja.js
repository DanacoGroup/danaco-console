/* Ekran trzeci kreatora instalacji umożliwia wybór wersji programu zgodnej z architekturą procesora wykrytą automatycznie przez instalator.

   Wykrycie procesora jest przesłanką, nie rozstrzygnięciem: instalator zaznacza
   wersję zgodną, ale zostawia wybór i ostrzega, gdy operator od wykrycia odchodzi.
   Trzy wyniki wykrycia (x64, arm, brak) niesie atrybut `data-wykryto` — ekran
   nie ma dla nich trzech gałęzi kodu. */
(function () {
'use strict';
var K = window.DanacoKreator, N = K.narzedzia, el = N.el, tekst = N.tekst, S = K.skladniki;
K.ekrany = K.ekrany || {};

K.ekrany[3] = function () {
  return el('section', {
    klasa: 'dn-kreator-ekran',
    dane: { ekran: 3, wykryto: 'x64' },
    'aria-label': tekst('szyna.kroki')[2]
  }, [
    S.naglowekEkranu({ nadtytul: 'krok3.nadtytul', tytul: 'krok3.tytul', podtytul: 'krok3.podtytulX64' }),
    S.tekstCiagly({ klucze: ['krok3.wskazowkaWykryta'], dane: { wskazowka: true } }),
    el('div', { klasa: 'dn-kreator-wybory', role: 'radiogroup', 'aria-label': tekst('krok3.wybor.etykieta') }, [
      S.kartaOpcji({ nazwa: 'krok3.x64.nazwa', opis: 'krok3.x64.opis', wartosc: 'x64', grupa: 'wariant', plakietka: 'krok3.plakietka', fokus: true }),
      S.kartaOpcji({ nazwa: 'krok3.arm.nazwa', opis: 'krok3.arm.opis', wartosc: 'arm', grupa: 'wariant', plakietka: 'krok3.plakietka' })
    ]),
    el('p', { klasa: 'dn-alert dn-alert--ostrzezenie', role: 'status', hidden: true, dane: { niezgodnosc: true }, tekst: tekst('krok3.niezgodnosc') }),
    S.odsylaczPomocy({ klucz: 'krok3.pomoc', dane: { 'pomoc-procesor': true } })
  ]);
};
})();
