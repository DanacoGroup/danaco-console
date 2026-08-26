/* ============================================================================
   SKŁADNIK — GŁOWA EKRANU

   Nadtytuł mówi, gdzie stoisz; tytuł — co się dzieje; zdanie wprowadzające —
   czego się spodziewać. Nadtytuł bywa zbędny: w oknie dostępu jego rolę pełnią
   zakładki, więc powtarzanie go byłoby szumem.

   Właściwości:
     nadtytul  klucz katalogu (opcjonalny)
     tytul     klucz katalogu
     lid       klucz katalogu (opcjonalny)
     lidWezly  wykaz węzłów zamiast samego tekstu — gdy w zdaniu stoi wyróżniony
               fragment, na przykład adres pisany krojem maszynowym
     poTytule  węzły wstawiane między tytuł a zdanie wprowadzające — tor kroków
               należy do głowy ekranu, nie stoi obok niej: głowa ma własny rytm
               odstępów i wyjęcie toru poza nią rozstraja go
     dane      atrybuty data-* na tytule
   ============================================================================ */
(function () {
'use strict';
var W = window.DanacoWejscie;
W.skladniki = W.skladniki || {};

W.skladniki.naglowekEkranu = function (N, w) {
  return N.el('div', { klasa: 'we-glowa' }, [
    w.nadtytul && N.el('p', { klasa: 'we-nadtytul', tekst: N.tekst(w.nadtytul) }),
    N.el('h2', { klasa: 'we-tytul', tekst: N.tekst(w.tytul), dane: w.dane || null })
  ].concat(w.poTytule || []).concat([
    w.lidWezly ? N.el('p', { klasa: 'we-lid' }, w.lidWezly)
      : (w.lid && N.el('p', { klasa: 'we-lid', tekst: w.lidTekst || N.tekst(w.lid) }))
  ]));
};
})();
