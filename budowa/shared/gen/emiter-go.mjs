// Emiter Go: sklada contract.go z deklaracji wytworzonych przez typy-go.mjs
// i blok-mapy-go.mjs. Wynik jest zgodny z gofmt bez uruchamiania gofmt.

import { linieNaglowka } from './naglowek.mjs';
import { mapaGo } from './blok-mapy-go.mjs';
import { opisDoBazy, opisZBazy } from './odwzorowanie-bazy.mjs';
import { narzedziaGo } from './narzedzia-go.mjs';
import {
  blokStalychGo,
  funkcjaListyGo,
  strukturaGo,
  typWyliczeniaGo,
} from './typy-go.mjs';

/** Komentarz pakietu wraz z ostrzezeniem o generacji oraz deklaracja pakietu. */
function naglowek(kontrakt) {
  const tresc = linieNaglowka(kontrakt, kontrakt.artefakty.go);
  return [
    `// Package ${kontrakt.artefakty.pakietGo} — kontrakt komunikacji ${kontrakt.produkt}.`,
    '//',
    ...tresc.map((l) => (l ? `// ${l}` : '//')),
    '',
    `package ${kontrakt.artefakty.pakietGo}`,
    '',
  ];
}

/** Import wyliczany z uzycia: surowy JSON i rozbior nazwy na obszar. */
function importy(model) {
  const pakiety = ['"strings"'];
  if (model.uzywaJson) pakiety.unshift('"encoding/json"');
  return ['import (', ...pakiety.map((p) => `\t${p}`), ')', ''];
}

/** Stale protokolu oraz typ nazwy komunikatu. */
function podstawy(model) {
  return [
    '// ProtocolVersion nazywa wersje protokolu kontraktu.',
    `const ProtocolVersion = "${model.kontrakt.protokol}"`,
    '',
    '// SeparatorObszaru oddziela obszar od reszty nazwy w notacji kontraktu.',
    `const SeparatorObszaru = "${model.kontrakt.separator}"`,
    '',
    '// MessageType nazywa komende albo zdarzenie. Jest aliasem napisu, wiec stale',
    '// kontraktu wchodza wprost tam, gdzie warstwa protokolu operuje na string.',
    'type MessageType = string',
    '',
  ];
}

/** Wyliczenia: nazwany typ napisowy i nietypowane stale wartosci. */
function wyliczenia(model) {
  return model.wyliczenia.flatMap((w) => [
    ...typWyliczeniaGo(w.nazwa, w.opis),
    ...blokStalychGo(`Wartosci ${w.nazwa}.`, w.wartosci.map((v) => ({ nazwa: v.staly, opis: v.opis, wartosc: v.wartosc }))),
    // Wykaz wartosci obok samych stalych. Rdzen sprawdzajacy wartosc przyslana
    // przez klienta musi ja z czyms porownac; bez tego wykazu kazde takie
    // sprawdzenie przepisywalo siedem napisow po raz kolejny i rozjezdzalo sie
    // z kontraktem przy pierwszej dolozonej wartosci.
    ...funkcjaListyGo(`Wartosci${w.nazwa}`, `zwraca komplet wartosci ${w.nazwa} w kolejnosci kontraktu.`,
      w.nazwa, w.wartosci.map((v) => v.staly)),
  ]);
}

/**
 * Wykaz "nazwa wyliczenia -> dopuszczalne wartosci" dla wszystkich wyliczen
 * kontraktu. Brama kontraktu rdzenia siega po niego po nazwie typu odczytanej
 * z pola zadania, wiec nowe wyliczenie wchodzi do bramy bez recznego wpisu.
 */
function zakresyWyliczen(model) {
  return mapaGo('ZakresyWyliczen', 'dopuszczalne wartosci kazdego wyliczenia kontraktu pod nazwa jego typu',
    'string', '[]string',
    model.wyliczenia.map((w) => [`"${w.nazwa}"`, `{${w.wartosci.map((v) => v.staly).join(', ')}}`]));
}

/**
 * Odwzorowanie wartosci wyliczenia na wartosc kolumny modelu danych i z powrotem.
 * Warstwa trwalosci siega po te slowniki zamiast wpisywac przeklad u siebie.
 */
function odwzorowania(model) {
  return model.odwzorowaniaBazy.flatMap((o) => [
    ...mapaGo(`WartosciBazy${o.nazwa}`, opisDoBazy(o), o.nazwa, 'string',
      o.wartosci.map((w) => [w.staly, `"${w.baza}"`])),
    ...mapaGo(`WartosciKontraktu${o.nazwa}`, opisZBazy(o), 'string', o.nazwa,
      o.wartosci.map((w) => [`"${w.baza}"`, w.staly])),
  ]);
}

/** Nazwy komend i zdarzen jako stale typu MessageType. */
function nazwyKomunikatow(model) {
  const komendy = model.komendy.map((k) => ({ nazwa: `Command${k.staly}`, opis: k.opis, wartosc: k.typ }));
  const zdarzenia = model.zdarzenia.map((z) => ({ nazwa: `Event${z.staly}`, opis: z.opis, wartosc: z.typ }));
  return [
    ...blokStalychGo('Nazwy komend kontraktu.', komendy, 'MessageType'),
    ...blokStalychGo('Nazwy zdarzen kontraktu.', zdarzenia, 'MessageType'),
  ];
}

/** Kody bledow i tablica ponawialnosci. */
function kodyBledow(model) {
  const wpisy = model.kodyBledow.map((b) => ({ nazwa: b.staly, opis: b.opis, wartosc: b.kod }));
  return [
    ...typWyliczeniaGo('ErrorCode', 'kod bledu odpowiedzi'),
    ...blokStalychGo('Kody bledow kontraktu.', wpisy),
    ...mapaGo(
      'KodyPonawialne',
      'czy ponowienie zadania po danym kodzie ma sens',
      'ErrorCode',
      'bool',
      model.kodyBledow.map((b) => [b.staly, String(b.retryable)]),
    ),
  ];
}

/** Koperta, struktury dziedzinowe i tresci komunikatow. */
function typyDanych(model) {
  const koperta = { nazwa: 'Envelope', opis: model.koperta.opis, pola: model.koperta.pola };
  return [
    ...strukturaGo(koperta),
    ...model.struktury.flatMap((s) => strukturaGo(s)),
    ...model.tresci.flatMap((t) => strukturaGo(t)),
  ];
}

/** Listy nazw, zbiory rozpoznania i zdarzenia zastepcze. */
function rozpoznanie(model) {
  const komendy = model.komendy.map((k) => `Command${k.staly}`);
  const zdarzenia = model.zdarzenia.map((z) => `Event${z.staly}`);
  return [
    ...funkcjaListyGo('WszystkieKomendy', 'zwraca komplet nazw komend kontraktu.', 'MessageType', komendy),
    ...funkcjaListyGo('WszystkieZdarzenia', 'zwraca komplet nazw zdarzen kontraktu.', 'MessageType', zdarzenia),
    ...mapaGo('zbiorKomend', 'Zbior nazw komend do rozpoznania typu.', 'MessageType', 'struct{}',
      komendy.map((k) => [k, '{}']), false),
    ...mapaGo('zbiorZdarzen', 'Zbior nazw zdarzen do rozpoznania typu.', 'MessageType', 'struct{}',
      zdarzenia.map((z) => [z, '{}']), false),
    ...mapaGo('zdarzeniaNieznanej', 'Zdarzenie *.unknown wlasciwe dla obszaru nazwy.', 'string', 'MessageType',
      model.zdarzenia.filter((z) => z.nieznane).map((z) => [`"${z.obszar}"`, `Event${z.staly}`]), false),
    '// CzyKomenda odpowiada, czy nazwa jest komenda kontraktu.',
    'func CzyKomenda(typ MessageType) bool {',
    '\t_, jest := zbiorKomend[typ]',
    '\treturn jest',
    '}',
    '',
    '// CzyZdarzenie odpowiada, czy nazwa jest zdarzeniem kontraktu.',
    'func CzyZdarzenie(typ MessageType) bool {',
    '\t_, jest := zbiorZdarzen[typ]',
    '\treturn jest',
    '}',
    '',
    '// ZdarzenieNieznanej zwraca zdarzenie zwracane na nieznana komende.',
    '// Fail-open: polaczenie nie jest zrywane, a obszar spoza',
    '// kontraktu dostaje zdarzenie zapasowe.',
    'func ZdarzenieNieznanej(typ MessageType) MessageType {',
    '\tobszar, _, _ := strings.Cut(typ, SeparatorObszaru)',
    '\tif zdarzenie, jest := zdarzeniaNieznanej[obszar]; jest {',
    '\t\treturn zdarzenie',
    '\t}',
    `\treturn Event${model.zdarzenieZapasowe}`,
    '}',
    '',
  ];
}

/** Listy wartosci znanych w chwili wydania kontraktu — informacyjne, nie bramy. */
function listyZnane(model) {
  return model.listyZnane.flatMap((l) =>
    funkcjaListyGo(l.nazwa, `zwraca liste: ${l.opis}.`, 'string', l.wartosci.map((w) => `"${w}"`)),
  );
}

export function emitujGo(model) {
  return [
    ...naglowek(model.kontrakt),
    ...importy(model),
    ...podstawy(model),
    ...wyliczenia(model),
    ...zakresyWyliczen(model),
    ...odwzorowania(model),
    ...nazwyKomunikatow(model),
    ...kodyBledow(model),
    ...typyDanych(model),
    ...rozpoznanie(model),
    ...narzedziaGo(model),
    ...listyZnane(model),
  ].join('\n');
}
