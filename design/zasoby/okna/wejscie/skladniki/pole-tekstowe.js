/* ============================================================================
   SKŁADNIK — POLE TEKSTOWE

   Etykieta stoi nad kontrolką, opis pod nią. Etykieta jest elementem `label`
   związanym identyfikatorem, nie samym tekstem obok — bez tego wskazanie
   etykiety nie ustawia kursora w polu, a czytnik ekranu nie wie, co czyta.

   Właściwości:
     etykieta   klucz katalogu
     id         identyfikator kontrolki
     typ        typ pola HTML (domyślnie 'text')
     uzupelnij  wartość autocomplete
     opis       klucz katalogu — zdanie pod polem
     bledne     true — kontrolka niesie `aria-invalid`; obwódka bierze się z tego
                stanu, nie z osobnej klasy
     opisuje    identyfikator komunikatu opisującego usterkę
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.poleTekstowe = function (N, w) {
  return N.el('label', { klasa: 'dn-pole' }, [
    N.el('span', { klasa: 'dn-pole-etykieta', tekst: N.tekst(w.etykieta) }),
    N.el('input', {
      klasa: 'dn-pole-kontrolka', type: w.typ || 'text', id: w.id,
      autocomplete: w.uzupelnij || null,
      'aria-invalid': w.bledne ? 'true' : null,
      'aria-describedby': w.opisuje || null
    }),
    w.opis && N.el('span', { klasa: 'dn-pole-opis', tekst: N.tekst(w.opis) })
  ]);
};
})();
