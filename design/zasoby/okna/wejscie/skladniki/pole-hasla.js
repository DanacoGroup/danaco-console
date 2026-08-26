/* ============================================================================
   SKŁADNIK — POLE HASŁA Z ODSŁONIĘCIEM

   Kontrolka odsłonięcia zmienia typ pola i własną etykietę, więc czytnik
   ekranu wie, w którym stanie stoi. Dwa znaki leżą w przycisku, a widoczność
   rozstrzyga arkusz po `aria-pressed` — przełączanie znaków skryptem
   rozjeżdżałoby się ze stanem kontrolki.

   Właściwości:
     etykieta   klucz katalogu
     id         identyfikator kontrolki
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

W.skladniki.poleHasla = function (N, w) {
  var odkryte = N.zeZnacznika(W.ikony.okoOdkryte);
  odkryte.setAttribute('class', 'ik-odkryte');
  odkryte.setAttribute('aria-hidden', 'true');
  var zakryte = N.zeZnacznika(W.ikony.okoZakryte);
  zakryte.setAttribute('class', 'ik-zakryte');
  zakryte.setAttribute('aria-hidden', 'true');

  return N.el('div', { klasa: 'dn-pole' }, [
    N.el('label', { klasa: 'dn-pole-etykieta', for: w.id, tekst: N.tekst(w.etykieta) }),
    N.el('span', { klasa: 'au-haslo' }, [
      N.el('input', {
        klasa: 'dn-pole-kontrolka', type: 'password', id: w.id,
        autocomplete: w.uzupelnij || 'current-password',
        'aria-invalid': w.bledne ? 'true' : null,
        'aria-describedby': w.opisuje || null
      }),
      N.el('button', {
        klasa: 'au-haslo-oko', type: 'button', 'aria-pressed': 'false',
        'aria-label': 'Pokaż hasło', dane: { odsloniecie: w.id }
      }, [odkryte, zakryte])
    ]),
    w.opis && N.el('span', { klasa: 'dn-pole-opis', id: w.id + '-opis', tekst: N.tekst(w.opis) })
  ]);
};
})();
