// Skrypt wnoszący pozycje kontraktu: wyłączenia pamięci, wyciszenie nakładki
// Always On Display wraz z nośnikiem sygnału klas zdarzeń oraz cztery brakujące
// pozycje modułu Studio.
//
// Skrypt jest IDEMPOTENTNY: pozycję o nazwie już obecnej pomija i mówi to
// wprost. Kopia `shared/contract.json` sprzed zmiany powstaje przed pierwszym
// zapisem, wzorem poprzednich dobudów.
import { readFileSync, writeFileSync, existsSync } from 'node:fs';

const SCIEZKA = '/home/ubuntu/DanacoConsole/budowa/shared/contract.json';
const KOPIA = SCIEZKA + '.przed-wyciszeniem-i-pamiecia';

const kontrakt = JSON.parse(readFileSync(SCIEZKA, 'utf8'));
if (!existsSync(KOPIA)) {
  writeFileSync(KOPIA, readFileSync(SCIEZKA));
  console.log('kopia sprzed zmiany: ' + KOPIA);
} else {
  console.log('kopia sprzed zmiany już stoi: ' + KOPIA);
}

const slad = [];
function odnotuj(co, nazwa, wniesione) {
  slad.push((wniesione ? 'wniesiono  ' : 'pominieto  ') + co + ' ' + nazwa);
}

// ── wyliczenia ───────────────────────────────────────────────────────────────

function wyliczenie(def) {
  const jest = kontrakt.wyliczenia.some((w) => w.nazwa === def.nazwa);
  if (!jest) kontrakt.wyliczenia.push(def);
  odnotuj('wyliczenie', def.nazwa, !jest);
}

function wartoscWyliczenia(nazwaWyliczenia, wartosc) {
  const w = kontrakt.wyliczenia.find((x) => x.nazwa === nazwaWyliczenia);
  if (!w) throw new Error('brak wyliczenia ' + nazwaWyliczenia);
  const jest = w.wartosci.some((x) => x.wartosc === wartosc.wartosc);
  if (!jest) w.wartosci.push(wartosc);
  odnotuj('wartosc', nazwaWyliczenia + '.' + wartosc.wartosc, !jest);
}

function struktura(def) {
  const jest = kontrakt.struktury.some((s) => s.nazwa === def.nazwa);
  if (!jest) kontrakt.struktury.push(def);
  odnotuj('struktura', def.nazwa, !jest);
}

function pole(nazwaStruktury, def, przed) {
  const s = kontrakt.struktury.find((x) => x.nazwa === nazwaStruktury);
  if (!s) throw new Error('brak struktury ' + nazwaStruktury);
  const jest = s.pola.some((p) => p.nazwa === def.nazwa);
  if (!jest) {
    const gdzie = przed === undefined ? -1 : s.pola.findIndex((p) => p.nazwa === przed);
    if (gdzie < 0) s.pola.push(def);
    else s.pola.splice(gdzie, 0, def);
  }
  odnotuj('pole', nazwaStruktury + '.' + def.nazwa, !jest);
}

function poleWyniku(typKomendy, def) {
  const k = kontrakt.komendy.find((x) => x.typ === typKomendy);
  if (!k) throw new Error('brak komendy ' + typKomendy);
  const jest = k.wynik.some((p) => p.nazwa === def.nazwa);
  if (!jest) k.wynik.push(def);
  odnotuj('pole wyniku', typKomendy + '.' + def.nazwa, !jest);
}

function komenda(def) {
  const jest = kontrakt.komendy.some((k) => k.typ === def.typ);
  if (!jest) kontrakt.komendy.push(def);
  odnotuj('komenda', def.typ, !jest);
}

function zdarzenie(def) {
  const jest = kontrakt.zdarzenia.some((z) => z.typ === def.typ);
  if (!jest) kontrakt.zdarzenia.push(def);
  odnotuj('zdarzenie', def.typ, !jest);
}

function opisKomendy(typ, opis) {
  const k = kontrakt.komendy.find((x) => x.typ === typ);
  if (!k) throw new Error('brak komendy ' + typ);
  const zmiana = k.opis !== opis;
  k.opis = opis;
  odnotuj('opis komendy', typ, zmiana);
}

const p = (nazwa, typ, wymagane, opis) => ({ nazwa, typ, wymagane, opis });

// ── Always On Display: klasy zdarzeń, rodzaj i zakres wyciszenia ─────────────

wyliczenie({
  nazwa: 'AodEventClass',
  opis:
    'Klasa zdarzenia wyzwalajacego sugestie nakladki — szesc wierszy tabeli rozdz. 3.2 ' +
    'opracowania always-on-display.md. Klasa jest jednoczesnie zakresem wyciszenia ' +
    '(rozdz. 3.5, wiersz trzeci): wyciszona klasa nie tworzy sugestii do chwili zniesienia, ' +
    'a pozostale dzialaja bez zmian',
  wartosci: [
    {
      wartosc: 'executionLoopState',
      opis:
        'Stan petli wykonawczej: petla wstrzymana, petla przerwana, zadanie ponawiane powyzej ' +
        'progu, zlecenie oczekujace na zatwierdzenie Uzytkownika',
    },
    {
      wartosc: 'taskQueueState',
      opis:
        'Stan kolejki zadan: kolejka zatrzymana, zadanie w stanie bledu, zadanie oczekujace ' +
        'dluzej niz prog czasu, kolejka wypelniona powyzej progu',
    },
    {
      wartosc: 'qualityControlResult',
      opis:
        'Wynik kontroli jakosci: negatywny wynik kontroli, powtorzone niepowodzenie tego samego ' +
        'zadania, rozbieznosc wyniku ze zleceniem',
    },
    {
      wartosc: 'moduleEvent',
      opis:
        'Zdarzenie modulu: zakonczenie dlugiego zadania, dostepny wynik do przegladu, konflikt ' +
        'zasobu, brak konfiguracji potrzebnej do wykonania polecenia',
    },
    {
      wartosc: 'schedule',
      opis:
        'Harmonogram: nadejscie terminu reguly czasowej, cykliczne uruchomienie automatyki, ' +
        'przypomnienie ustanowione przez Operatora',
    },
    {
      wartosc: 'operatorWorkContext',
      opis:
        'Kontekst pracy Operatora: powtarzalnosc recznie wykonywanej czynnosci, praca w module ' +
        'bez skonfigurowanego wykonawcy, otwarty zasob powiazany z zadaniem oczekujacym',
    },
  ],
});

wyliczenie({
  nazwa: 'AodMuteKind',
  opis:
    'Rodzaj wyciszenia nakladki — trzy pierwsze wiersze tabeli rozdz. 3.5 opracowania ' +
    'always-on-display.md. Tryb cichy NIE JEST czwarta wartoscia tego wyliczenia: jest trybem ' +
    'obecnosci (rozdz. 9.2) i jedzie ustawieniem konfiguracji, nie wpisem wyciszenia. Wyjatek ' +
    'wagi krytycznej takze nie jest wyciszeniem — jest regula przebijajaca kazde z trzech',
  wartosci: [
    {
      wartosc: 'timed',
      opis:
        'Wyciszenie czasowe: sugestie gromadza sie w liscie oczekujacych, dymek sie nie otwiera, ' +
        'plakietka pozostaje ukryta. Chwile konca niesie pole endsAt',
    },
    {
      wartosc: 'contextual',
      opis:
        'Wyciszenie kontekstowe: sugestie wskazanego modulu albo wskazanej karty sesji nie ' +
        'ujawniaja sie, pozostale zachowuja pelne dzialanie',
    },
    {
      wartosc: 'eventClass',
      opis:
        'Wyciszenie klasy zdarzen: wskazana klasa nie tworzy sugestii do chwili zniesienia, ' +
        'pozostale klasy dzialaja bez zmian',
    },
  ],
});

wyliczenie({
  nazwa: 'AodMuteScope',
  opis:
    'Zakres wyciszenia nakladki — kolumna „Zakres" tabeli rozdz. 3.5. Wyciszenie czasowe ' +
    'obowiazuje cala platforme i zakresu nie niesie (pole scope puste); wyciszenie kontekstowe ' +
    'wskazuje modul albo karte sesji, wyciszenie klasy zdarzen — klase',
  wartosci: [
    { wartosc: 'module', opis: 'Biezacy modul; byt wskazuje scopeId' },
    { wartosc: 'session', opis: 'Biezaca karta sesji; byt wskazuje scopeId' },
    { wartosc: 'eventClass', opis: 'Klasa zdarzen wskazana polem eventClass' },
  ],
});

struktura({
  nazwa: 'AodMute',
  opis:
    'Jedno wyciszenie nakladki czynne w rdzeniu. Wyciszenie jest bytem rdzenia, nie stanem ' +
    'jednego okna: Operator wyciszajacy nakladke w jednej powloce ma miec cisze w kazdej ' +
    'nastepnej. Zniesienie idzie ta sama komenda aod.mute.set z polem muted rownym falszowi ' +
    '— jednym ruchem, bez pytania o potwierdzenie',
  pola: [
    p('id', 'string', true, 'Identyfikator wyciszenia; podaje go zniesienie'),
    p('kind', 'AodMuteKind', true, 'Rodzaj wyciszenia'),
    p(
      'scope',
      'AodMuteScope',
      false,
      'Zakres wyciszenia; puste dla wyciszenia czasowego, ktore obowiazuje cala platforme',
    ),
    p('scopeId', 'string', false, 'Modul albo karta sesji objeta wyciszeniem kontekstowym'),
    p('scopeName', 'string', false,
      'Nazwa bytu zakresu widziana przez Operatora — pelna nazwa modulu albo karty sesji, nie kod. ' +
      'Cisza, przy ktorej Operator nie wie, CO milczy, jest gorsza od braku wyciszenia'),
    p('eventClass', 'AodEventClass', false, 'Klasa zdarzen objeta wyciszeniem klasy'),
    p(
      'endsAt',
      'int64',
      false,
      'Chwila konca wyciszenia w milisekundach epoki. Puste znaczy wyciszenie trwajace do ' +
        'zniesienia reka Operatora — tak stoi wyciszenie kontekstowe i wyciszenie klasy zdarzen',
    ),
    p(
      'deviceId',
      'string',
      false,
      'Urzadzenie, z ktorego wyciszenie zalozono. Wyciszenie obowiazuje wszystkie powloki ' +
        'Operatora niezaleznie od tego pola; pole mowi, skad przyszlo',
    ),
    p('createdAt', 'int64', true, 'Czas zalozenia w milisekundach epoki'),
  ],
});

struktura({
  nazwa: 'AodSignal',
  opis:
    'Sygnal klasy zdarzen wyzwalajacych, odlozony w rdzeniu. Nosnik dla klas, ktorych rdzen nie ' +
    'widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow automatyk ' +
    'i powtarzalnosci czynnosci Operatora (rozdz. 3.2 i 3.3 opracowania). Bez niego trzy z ' +
    'szesciu klas nie maja czego wyciszac ani z czego zrobic sugestii',
  pola: [
    p('id', 'string', true, 'Identyfikator sygnalu'),
    p('eventClass', 'AodEventClass', true, 'Klasa zdarzenia wyzwalajacego'),
    p('text', 'string', true, 'Tresc sygnalu widziana przez Operatora'),
    p('moduleId', 'string', false, 'Modul, ktorego sygnal dotyczy'),
    p('sessionId', 'string', false, 'Karta sesji, ktorej sygnal dotyczy'),
    p(
      'occurrenceCount',
      'int',
      false,
      'Liczba wystapien tej samej czynnosci. Prog powtarzalnosci czynnosci recznej z rozdz. 3.4 ' +
        'liczy sie z tego pola; puste znaczy jedno wystapienie',
    ),
    p('occurredAt', 'int64', true, 'Czas zdarzenia w milisekundach epoki'),
  ],
});

pole('AodStatus', p(
  'moduleId',
  'string',
  false,
  'Modul pokazywany w nakladce. Do tej pozycji nakladka dochodzila modulu okreznie — przez ' +
    'window.list — i sugestia, ktorej okna rdzen nie znal, w wyciszenie modulu nie wpadala',
), 'attachedProcessIds');

pole('AodSuggestion', p(
  'moduleId',
  'string',
  false,
  'Modul, ktorego podpowiedz dotyczy. Bez niego wyciszenie biezacego modulu nie ma po czym ' +
    'rozpoznac swojej sugestii',
), 'createdAt');

pole('AodSuggestion', p(
  'eventClass',
  'AodEventClass',
  false,
  'Klasa zdarzenia, ktore podpowiedz wywolalo. Podstawa wyciszenia klasy zdarzen; puste znaczy ' +
    'klase nierozpoznana, a sugestia bez klasy w wyciszenie klasy nie wpada',
), 'createdAt');

komenda({
  typ: 'aod.mute.get',
  opis:
    'Zwraca wyciszenia nakladki czynne w rdzeniu wraz z rodzajem, zakresem i chwila konca. ' +
    'Wyciszenie przeterminowane nie wchodzi do wykazu — Operator nie ma go odklikiwac',
  zadanie: [
    p('deviceId', 'string', false, 'Urzadzenie pytajace; wyciszenia obowiazuja wszystkie powloki'),
  ],
  wynik: [
    p('mutes', 'AodMute[]', true, 'Wyciszenia czynne w kolejnosci zalozenia'),
  ],
});

komenda({
  typ: 'aod.mute.set',
  opis:
    'Zaklada wyciszenie nakladki albo je znosi. Zniesienie idzie ta sama komenda z polem muted ' +
    'rownym falszowi — odwracalnosc jednym ruchem jest wymogiem, nie wygoda. Odpowiedz oddaje ' +
    'wykaz wyciszen PO zmianie, zeby powloka nie musiala pytac drugi raz. Zmiana rozglasza sie ' +
    'zdarzeniem aod.mute.changed na pozostale powloki Operatora',
  zadanie: [
    p(
      'muteId',
      'string',
      false,
      'Wyciszenie znoszone. Podane przy muted rownym falszowi znosi dokladnie to wyciszenie; ' +
        'puste znosi wyciszenie zlozone z rodzaju i zakresu tego zadania',
    ),
    p('kind', 'AodMuteKind', false, 'Rodzaj wyciszenia; wymagany przy zalozeniu'),
    p('scope', 'AodMuteScope', false, 'Zakres wyciszenia; wymagany przy wyciszeniu kontekstowym'),
    p('scopeId', 'string', false, 'Modul albo karta sesji objeta wyciszeniem kontekstowym'),
    p('scopeName', 'string', false, 'Nazwa bytu zakresu widziana przez Operatora'),
    p('eventClass', 'AodEventClass', false, 'Klasa zdarzen; wymagana przy wyciszeniu klasy'),
    p(
      'endsAt',
      'int64',
      false,
      'Chwila konca wyciszenia czasowego w milisekundach epoki. Wyciszenie czasowe bez niej jest ' +
        'zadaniem niezgodnym z kontraktem — cisza bez konca nie mowi Operatorowi, do kiedy trwa',
    ),
    p('muted', 'bool', true, 'Prawda zaklada wyciszenie, falsz je znosi'),
    p('deviceId', 'string', false, 'Urzadzenie, z ktorego idzie zmiana'),
  ],
  wynik: [
    p('mutes', 'AodMute[]', true, 'Wyciszenia czynne po zmianie'),
    p(
      'changed',
      'bool',
      true,
      'Czy wykaz naprawde sie zmienil. Falsz znaczy wyciszenie powtorzone albo zniesienie ' +
        'wyciszenia, ktorego nie bylo',
    ),
  ],
});

komenda({
  typ: 'aod.signal.report',
  opis:
    'Odklada w rdzeniu sygnal klasy zdarzen wyzwalajacych. Nosnik dla trzech klas rozdz. 3.2, ' +
    'ktorych rdzen nie widzi wlasna telemetria: wyniku kontroli jakosci, harmonogramu przebiegow ' +
    'automatyk i powtarzalnosci czynnosci Operatora. Odpowiedz mowi wprost, czy sygnal wpadl ' +
    'w wyciszenie i ktore — sygnal wyciszony odklada sie nadal, bo wyciszenie wstrzymuje ' +
    'UJAWNIENIE, a nie zapis',
  zadanie: [
    p('eventClass', 'AodEventClass', true, 'Klasa zdarzenia wyzwalajacego'),
    p('text', 'string', true, 'Tresc sygnalu widziana przez Operatora'),
    p('moduleId', 'string', false, 'Modul, ktorego sygnal dotyczy'),
    p('sessionId', 'string', false, 'Karta sesji, ktorej sygnal dotyczy'),
    p('occurrenceCount', 'int', false, 'Liczba wystapien tej samej czynnosci; puste znaczy jedno'),
  ],
  wynik: [
    p('signal', 'AodSignal', true, 'Sygnal po odlozeniu'),
    p('suppressed', 'bool', true, 'Czy sygnal wpadl w wyciszenie i nie ujawni sie w nakladce'),
    p('mute', 'AodMute', false, 'Wyciszenie, ktore sygnal wstrzymalo; puste, gdy zadne'),
  ],
});

komenda({
  typ: 'aod.signal.list',
  opis:
    'Zwraca sygnaly klas zdarzen odlozone w rdzeniu, bez tych wstrzymanych wyciszeniem. ' +
    'Odpowiedz nazywa wyciszenia, ktore sygnaly wstrzymaly, i liczbe wstrzymanych — cisza, po ' +
    'ktorej Operator nie wie, ze cos milczy, jest gorsza od braku wyciszenia',
  zadanie: [
    p('deviceId', 'string', false, 'Urzadzenie pytajace'),
    p('classes', 'AodEventClass[]', false, 'Zawezenie do wskazanych klas; puste znaczy wszystkie'),
    p('moduleId', 'string', false, 'Zawezenie do jednego modulu'),
    p('sessionId', 'string', false, 'Zawezenie do jednej karty sesji'),
    p('limit', 'int', false, 'Gorna granica liczby sygnalow'),
  ],
  wynik: [
    p('signals', 'AodSignal[]', true, 'Sygnaly ujawnialne, najswiezsze na poczatku'),
    p(
      'suppressedCount',
      'int',
      true,
      'Liczba sygnalow wstrzymanych wyciszeniem, nazwana zamiast przemilczana',
    ),
    p('mutes', 'AodMute[]', true, 'Wyciszenia, ktore sygnaly wstrzymaly'),
  ],
});

zdarzenie({
  typ: 'aod.mute.changed',
  opis:
    'Wyciszenie nakladki zalozone albo zniesione. Rozgloszenie sciga cisza pozostale powloki ' +
    'Operatora: bez niego Operator wyciszalby w jednej, a sugestie wchodzilyby w drugiej',
  payload: [
    p('change', 'ChangeKind', true, 'Rodzaj zmiany'),
    p('mute', 'AodMute', true, 'Wyciszenie objete zmiana'),
    p('mutes', 'AodMute[]', true, 'Wyciszenia czynne po zmianie'),
  ],
});

// ── Pamięć: wyłączenia ──────────────────────────────────────────────────────

struktura({
  nazwa: 'MemoryDisable',
  opis:
    'Wylaczenie pamieci w zasiegu. Wylacza sie albo pojedynczy wpis (entryId), albo caly poziom ' +
    'pamieci (level) — w zasiegu: calkiem, w srodowisku, w projekcie, w module, w parze modulow ' +
    'albo w karcie sesji. WYLACZENIE NIE KASUJE TRESCI i to jest jego roznica wobec memory.delete: ' +
    'wpis wylaczony nie wchodzi do kontekstu, ale zostaje na miejscu i wraca w calosci ' +
    'zniesieniem wylaczenia. Sterowanie stoi po stronie Operatora w konfiguracji — zakres ' +
    'ustawien „Pamiec"',
  pola: [
    p('id', 'string', true, 'Identyfikator wylaczenia; podaje go zniesienie'),
    p(
      'entryId',
      'string',
      false,
      'Wpis pamieci objety wylaczeniem. Puste znaczy wylaczenie calego poziomu wskazanego polem level',
    ),
    p(
      'level',
      'MemoryLevel',
      false,
      'Poziom pamieci objety wylaczeniem. Puste znaczy wylaczenie pojedynczego wpisu wskazanego ' +
        'polem entryId. Zadanie bez obu pol nie ma czego wylaczyc i konczy sie odmowa',
    ),
    p(
      'scope',
      'ConfigScope',
      true,
      'Zasieg wylaczenia: global znaczy „calkiem", environment — w srodowisku, project ' +
        '— w projekcie, module — w module, modulePair — w parze modulow, session — w karcie sesji',
    ),
    p('scopeId', 'string', false, 'Byt zasiegu; puste dla zasiegu global'),
    p('createdAt', 'int64', true, 'Czas zalozenia w milisekundach epoki'),
    p('updatedAt', 'int64', true, 'Czas ostatniej zmiany w milisekundach epoki'),
  ],
});

struktura({
  nazwa: 'MemoryDisabledEntry',
  opis:
    'Wpis pamieci wstrzymany wylaczeniem: nazwany wprost obok wykazu wpisow czynnych, wraz ' +
    'z zasiegiem, ktory go wylaczyl. Wpis znika z kontekstu, nie z pamieci — tresc stoi ' +
    'nietknieta i wraca zniesieniem wylaczenia',
  pola: [
    p('entryId', 'string', true, 'Wpis pamieci, ktory nie wszedl do kontekstu'),
    p('disableId', 'string', true, 'Wylaczenie, ktore wpis wstrzymalo; podaje je zniesienie'),
    p('scope', 'ConfigScope', true, 'Zasieg, ktory wpis wylaczyl'),
    p('scopeId', 'string', false, 'Byt zasiegu; puste dla zasiegu global'),
    p('level', 'MemoryLevel', false, 'Poziom pamieci, gdy wylaczono caly poziom, a nie sam wpis'),
  ],
});

poleWyniku('memory.list', p(
  'disabledEntries',
  'MemoryDisabledEntry[]',
  true,
  'Wpisy wstrzymane wylaczeniem, wraz z zasiegiem, ktory je wylaczyl. Do kontekstu nie wchodza ' +
    'i nie stoja w polu entries; wykaz jest zawsze tablica, bo brak wylaczen to fakt, nie cisza',
));

komenda({
  typ: 'memory.disable.list',
  opis:
    'Zwraca wylaczenia pamieci obowiazujace w zasiegu zadania. Odczyt do pary z ' +
    'memory.disable.set; okno konfiguracji buduje z niego zakres ustawien „Pamiec"',
  zadanie: [
    p('projectId', 'string', false, 'Projekt, ktorego wylaczenia sa czytane'),
    p('sessionId', 'string', false, 'Karta sesji, gdy projekt nie zostal wskazany'),
    p('entryId', 'string', false, 'Zawezenie do wylaczen jednego wpisu'),
    p('level', 'MemoryLevel', false, 'Zawezenie do wylaczen jednego poziomu pamieci'),
    p('scope', 'ConfigScope', false, 'Zawezenie do jednego zasiegu wylaczenia'),
    p('scopeId', 'string', false, 'Byt zasiegu; puste dla zasiegu global'),
  ],
  wynik: [
    p('disables', 'MemoryDisable[]', true, 'Wylaczenia obowiazujace w zasiegu zadania'),
  ],
});

komenda({
  typ: 'memory.disable.set',
  opis:
    'Wylacza wpis pamieci albo caly poziom pamieci w zasiegu — calkiem, w srodowisku, ' +
    'w projekcie, w module, w parze modulow albo w karcie sesji — albo wylaczenie znosi. ' +
    'Zniesienie idzie ta sama komenda z polem disabled rownym falszowi: odwracalnosc jednym ' +
    'ruchem. WYLACZENIE NIE KASUJE TRESCI — tym rozni sie od memory.delete; wpis wylaczony ' +
    'zostaje na miejscu i wraca w calosci po zniesieniu. Odpowiedz oddaje wykaz wylaczen PO zmianie',
  zadanie: [
    p(
      'disableId',
      'string',
      false,
      'Wylaczenie znoszone. Puste przy disabled rownym falszowi znosi wylaczenie zlozone z pol ' +
        'tego zadania',
    ),
    p('entryId', 'string', false, 'Wpis pamieci wylaczany; puste znaczy wylaczenie calego poziomu'),
    p('level', 'MemoryLevel', false, 'Poziom pamieci wylaczany; puste znaczy wylaczenie wpisu'),
    p('scope', 'ConfigScope', true, 'Zasieg wylaczenia'),
    p('scopeId', 'string', false, 'Byt zasiegu; puste dla zasiegu global'),
    p('disabled', 'bool', true, 'Prawda zaklada wylaczenie, falsz je znosi'),
  ],
  wynik: [
    p('disable', 'MemoryDisable', true, 'Wylaczenie po zmianie'),
    p('disables', 'MemoryDisable[]', true, 'Wylaczenia obowiazujace po zmianie'),
    p(
      'changed',
      'bool',
      true,
      'Czy wykaz naprawde sie zmienil. Falsz znaczy wylaczenie powtorzone albo zniesienie ' +
        'wylaczenia, ktorego nie bylo',
    ),
  ],
});

zdarzenie({
  typ: 'memory.disable.changed',
  opis:
    'Wylaczenie pamieci zalozone albo zniesione. Rozgloszenie idzie na pozostale powloki ' +
    'Operatora, zeby wykaz pamieci w kazdej z nich pokazywal ten sam stan wlaczenia',
  payload: [
    p('change', 'ChangeKind', true, 'Rodzaj zmiany'),
    p('disable', 'MemoryDisable', true, 'Wylaczenie objete zmiana'),
  ],
});

opisKomendy(
  'memory.detach',
  'Odpina wpis pamieci od zasiegu wspoldzielonego, sprowadzajac go do projektu, ktory go niesie; ' +
    'tresc zostaje nietknieta, a memory.set zasieg przywraca. Odpiecie NIE JEST wylaczeniem: ' +
    'zweza zasieg samego wpisu, a nie wstrzymuje go w cudzym projekcie ani module — do tego ' +
    'sluzy memory.disable.set',
);

// ── Studio: cztery brakujące pozycje ────────────────────────────────────────

pole('StudioFragmentLock', p(
  'templateId',
  'string',
  false,
  'Szablon, z ktorego blokada wzorcowa pochodzi. Blokada wychodzaca z fromTemplate bez tego pola ' +
    'nie mowi, KTOREGO szablonu jest wzorcem, wiec Operator nie ma jak dojsc do jej zrodla',
));

pole('StudioParagraphFormat', p(
  'listRestart',
  'bool',
  false,
  'Czy numeracja listy wznawia sie od tego akapitu. Bez tego pola wznowienie rozdzielalo liste na ' +
    'dwie definicje, a zmiana jednej definicji przestawala przestawiac oba jej odcinki',
));

wartoscWyliczenia('StudioActionKind', {
  wartosc: 'fieldChange',
  opis:
    'Zmiana pola dokumentu. Osobno od objectChange, bo pole odlozone jako zmiana obiektu nie da ' +
    'sie cofnac samo — Operator cofalby wtedy cala zmiane obiektu razem z nia',
});

// ── zapis ───────────────────────────────────────────────────────────────────

writeFileSync(SCIEZKA, JSON.stringify(kontrakt, null, 2) + '\n');
for (const wiersz of slad) console.log(wiersz);
console.log(
  'komend: ' + kontrakt.komendy.length + ', zdarzen: ' + kontrakt.zdarzenia.length +
    ', struktur: ' + kontrakt.struktury.length + ', wyliczen: ' + kontrakt.wyliczenia.length,
);
