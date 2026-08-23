// Odwzorowanie wartosci kontraktu na wartosci kolumn modelu danych.
// Kontrakt zapisuje nazwy po angielsku, baza po polsku; przeklad ma byc jeden
// do jednego i zyc w jednym miejscu — w sekcji `wyliczenia` pliku contract.json.
// Modul nie zna skladni Go ani TypeScript: sprawdza kompletnosc i wzajemna
// jednoznacznosc odwzorowania oraz podaje emiterom gotowe pary.

/** Wyliczenie ma odwzorowanie, jezeli wskazuje kolumne modelu danych. */
function maOdwzorowanie(wyliczenie) {
  return typeof wyliczenie.kolumnaBazy === 'string' && wyliczenie.kolumnaBazy.length > 0;
}

/**
 * Wartosc przelotowa nie ma odpowiednika w kolumnie modelu danych i tak ma byc.
 * Sluzy wylacznie przesylowi (np. fragment diagnostyczny strumienia), wiec nie
 * moze trafic do kolumny, ktorej CHECK zna tylko wartosci utrwalane. Znacznik
 * jest jawny — brak pola `baza` bez `przelotowa` nadal jest bledem kontraktu.
 */
function przelotowa(w) {
  return w.przelotowa === true;
}

/** Kazda wartosc utrwalana musi miec odpowiednik w bazie — inaczej przeklad jest dziurawy. */
function sprawdzKompletnosc({ nazwa, wartosci }) {
  for (const w of wartosci) {
    if (przelotowa(w)) {
      if (typeof w.baza === 'string' && w.baza.length > 0) {
        throw new Error(
          `contract.json: wyliczenie ${nazwa}, wartosc "${w.wartosc}" jest przelotowa, wiec nie moze miec pola "baza"`,
        );
      }
      continue;
    }
    if (typeof w.baza !== 'string' || w.baza.length === 0) {
      throw new Error(
        `contract.json: wyliczenie ${nazwa} wskazuje kolumne bazy, ale wartosc "${w.wartosc}" nie ma pola "baza"`,
      );
    }
  }
}

/** Dwie wartosci kontraktu nie moga wskazywac tej samej wartosci w bazie — odwrotnosc byla by niejednoznaczna. */
function sprawdzJednoznacznosc({ nazwa, wartosci }) {
  const zajete = new Map();
  for (const w of wartosci) {
    if (zajete.has(w.baza)) {
      throw new Error(
        `contract.json: wyliczenie ${nazwa} odwzorowuje "${zajete.get(w.baza)}" i "${w.wartosc}" na te sama wartosc bazy "${w.baza}"`,
      );
    }
    zajete.set(w.baza, w.wartosc);
  }
}

/**
 * Opis odwzorowania jednego wyliczenia: nazwa typu, wskazana kolumna oraz pary
 * gotowe do zapisania w obu kierunkach.
 */
function opisOdwzorowania(wyliczenie) {
  sprawdzKompletnosc(wyliczenie);
  const utrwalane = wyliczenie.wartosci.filter((w) => !przelotowa(w));
  sprawdzJednoznacznosc({ nazwa: wyliczenie.nazwa, wartosci: utrwalane });
  return {
    nazwa: wyliczenie.nazwa,
    kolumna: wyliczenie.kolumnaBazy,
    // Pelne odwzorowanie pokrywa kazda wartosc wyliczenia. Wartosc przelotowa
    // odpowiednika w kolumnie nie ma, wiec slownik przestaje byc pelny.
    pelne: utrwalane.length === wyliczenie.wartosci.length,
    wartosci: utrwalane.map((w) => ({
      wartosc: w.wartosc,
      baza: w.baza,
      staly: w.staly,
    })),
  };
}

/** Odwzorowania wszystkich wyliczen majacych odpowiednik w modelu danych. */
export function odwzorowaniaBazy(wyliczenia) {
  return wyliczenia.filter(maOdwzorowanie).map(opisOdwzorowania);
}

/** Komentarz kierunku kontrakt -> baza. */
export function opisDoBazy({ nazwa, kolumna }) {
  return `wartosc kontraktu ${nazwa} -> wartosc kolumny ${kolumna}`;
}

/** Komentarz kierunku baza -> kontrakt. */
export function opisZBazy({ nazwa, kolumna }) {
  return `wartosc kolumny ${kolumna} -> wartosc kontraktu ${nazwa}`;
}
