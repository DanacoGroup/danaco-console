// Model kontraktu: normalizacja pliku contract.json do postaci, ktora emitery
// renderuja wprost, bez wlasnych decyzji. Tu rozwijane sa zdarzenia *.unknown.

import { odwzorowaniaBazy } from './odwzorowanie-bazy.mjs';
import { narzedziaModelu } from './narzedzia.mjs';
import {
  naPascal,
  rozbijTyp,
  nazwaZadania,
  nazwaWyniku,
  nazwaZdarzenia,
  stalaWyliczenia,
  NAZWA_TRESCI_NIEZNANEJ,
} from './nazwy.mjs';

/** Koperta jest jedna: pola wspolne, odpowiedzi i strumienia w jednym typie. */
function polaKoperty(koperta) {
  return [
    ...koperta.polaWspolne,
    ...koperta.polaOdpowiedzi,
    ...koperta.polaStrumienia,
  ];
}

/**
 * Wyliczenia z nazwami stalych dla kazdej wartosci; `baza` i znacznik
 * `przelotowa` ida dalej bez zmian — o odwzorowaniu na kolumne rozstrzyga
 * modul odwzorowania, nie ten.
 */
function wyliczenia(zrodlo) {
  return zrodlo.wyliczenia.map((w) => ({
    nazwa: w.nazwa,
    opis: w.opis,
    kolumnaBazy: w.kolumnaBazy,
    wartosci: w.wartosci.map((v) => ({
      wartosc: v.wartosc,
      baza: v.baza,
      przelotowa: v.przelotowa,
      opis: v.opis,
      staly: stalaWyliczenia(w.nazwa, v.wartosc),
    })),
  }));
}

/** Tresci komunikatow: zadanie i wynik kazdej komendy, tresc kazdego zdarzenia. */
function tresci(zrodlo) {
  const lista = [];
  for (const k of zrodlo.komendy) {
    lista.push({ nazwa: nazwaZadania(k.typ), opis: `Tresc zadania ${k.typ} — ${k.opis}`, pola: k.zadanie });
    lista.push({ nazwa: nazwaWyniku(k.typ), opis: `Tresc wyniku ${k.typ} — ${k.opis}`, pola: k.wynik });
  }
  for (const z of zrodlo.zdarzenia) {
    lista.push({ nazwa: nazwaZdarzenia(z.typ), opis: `Tresc zdarzenia ${z.typ} — ${z.opis}`, pola: z.payload });
  }
  lista.push({
    nazwa: NAZWA_TRESCI_NIEZNANEJ,
    opis: zrodlo.zdarzeniaNieznane.opis,
    pola: zrodlo.zdarzeniaNieznane.payload,
  });
  return lista;
}

/** Komendy z nazwa stalej i nazwami typow tresci. */
function komendy(zrodlo) {
  return zrodlo.komendy.map((k) => ({
    typ: k.typ,
    opis: k.opis,
    staly: naPascal(k.typ),
    typZadania: nazwaZadania(k.typ),
    typWyniku: nazwaWyniku(k.typ),
    polaWymagane: (k.zadanie ?? []).filter((p) => p.wymagane).map((p) => p.nazwa),
  }));
}

/** Zdarzenia zadeklarowane wprost oraz rozwiniete *.unknown dla kazdego obszaru. */
function zdarzenia(zrodlo) {
  const jawne = zrodlo.zdarzenia.map((z) => ({
    typ: z.typ,
    opis: z.opis,
    staly: naPascal(z.typ),
    typTresci: nazwaZdarzenia(z.typ),
    obszar: z.typ.split('.')[0],
    nieznane: false,
  }));
  const { przyrostek, opis } = zrodlo.zdarzeniaNieznane;
  const nieznane = zrodlo.obszary.map((o) => ({
    typ: `${o.nazwa}.${przyrostek}`,
    opis: `${o.opis} — ${opis}`,
    staly: naPascal(`${o.nazwa}.${przyrostek}`),
    typTresci: NAZWA_TRESCI_NIEZNANEJ,
    obszar: o.nazwa,
    nieznane: true,
  }));
  return [...jawne, ...nieznane];
}

/** Zbior nazw uzytkownikowych typow zlozonych — potrzebny emiterowi Go do wskaznikow. */
function nazwyStruktur(zrodlo, listaTresci) {
  return new Set([...zrodlo.struktury.map((s) => s.nazwa), ...listaTresci.map((t) => t.nazwa)]);
}

/** Czy w kontrakcie wystepuje pole o typie surowego JSON. */
function uzywaJson(zbioryPol) {
  return zbioryPol.some((pola) => pola.some((p) => rozbijTyp(p.typ).bazowy === 'json'));
}

export function zbudujModel(zrodlo) {
  const listaTresci = tresci(zrodlo);
  const koperta = polaKoperty(zrodlo.koperta);
  const wszystkiePola = [koperta, ...zrodlo.struktury.map((s) => s.pola), ...listaTresci.map((t) => t.pola)];
  const obszarZapasowy = zrodlo.zdarzeniaNieznane.obszarZapasowy;
  const listaWyliczen = wyliczenia(zrodlo);
  return {
    kontrakt: zrodlo.kontrakt,
    koperta: { opis: zrodlo.koperta.opis, pola: koperta },
    wyliczenia: listaWyliczen,
    odwzorowaniaBazy: odwzorowaniaBazy(listaWyliczen),
    struktury: zrodlo.struktury,
    tresci: listaTresci,
    komendy: komendy(zrodlo),
    zdarzenia: zdarzenia(zrodlo),
    narzedzia: narzedziaModelu(zrodlo),
    obszary: zrodlo.obszary,
    zdarzenieZapasowe: naPascal(`${obszarZapasowy}.${zrodlo.zdarzeniaNieznane.przyrostek}`),
    kodyBledow: zrodlo.kodyBledow.map((b) => ({ ...b, staly: `ErrorCode${naPascal(b.kod)}` })),
    listyZnane: zrodlo.listyZnane,
    nazwyStruktur: nazwyStruktur(zrodlo, listaTresci),
    nazwyWyliczen: new Set(zrodlo.wyliczenia.map((w) => w.nazwa)),
    uzywaJson: uzywaJson(wszystkiePola),
  };
}
