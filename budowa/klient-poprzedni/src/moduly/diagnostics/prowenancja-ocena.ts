import { ModelCallQuality, type ModelCallTrace } from '../../../../shared/contract';
import { poleTresci, przyciskAkcji, wiersz, wybor } from '../../modele/kontrolki-formularza-braki';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import { OCENA_WYWOLANIA } from './prowenancja-slowa';
import type { ZrodloProwenancji } from './prowenancja-zrodlo';

/**
 * Ocena wywołania modelu — czynność Operatora, nie odczyt.
 *
 * Rodzina prowenancji ma cztery komendy odczytu i jedną, która zapisuje sąd
 * człowieka o pracy modelu: `provenance.call.rate`. Dlatego ocena ma w oknie
 * własny, jawny chwyt — pole wyboru trafności, pole uzasadnienia i przycisk
 * zapisu — a nie skrót klawiszowy ani kliknięcie w wiersz. Czynność, która
 * zapisuje ocenę po naciśnięciu czegoś, co wygląda jak odczyt, jest czynnością
 * ukrytą.
 *
 * Opracowanie modułu mówi „pole oceny” i „adnotacja jakości”, ale skali nie
 * ustala. Skala pochodzi więc z kontraktu (`ModelCallQuality`) i to jest jedyne
 * miejsce, w którym wolno ją wziąć: skala wymyślona w oknie nie miałaby gdzie
 * się zapisać. Wartość „bez oceny” zostaje w wyborze celowo — Operator, który
 * ocenił omyłkowo, musi mieć drogę zdjęcia oceny, a kontrakt tę wartość niesie.
 *
 * Uzasadnienie jest nieobowiązkowe w kontrakcie (`note`) i takie zostaje tutaj:
 * wymuszenie go w oknie byłoby zaporą, której rdzeń nie stawia, a ocena bez
 * słowa nadal jest oceną.
 */
export interface OcenaWywolania {
  element: HTMLElement;
}

/** Trafności w kolejności od najlepszej; „bez oceny” stoi na końcu jako zdjęcie oceny. */
const SKALA: readonly ModelCallQuality[] = [
  ModelCallQuality.Accurate,
  ModelCallQuality.Partial,
  ModelCallQuality.Inaccurate,
  ModelCallQuality.Unrated,
];

export function utworzOceneWywolania(
  zrodlo: ZrodloProwenancji,
  wywolanie: ModelCallTrace,
  poZapisie: (ocenione: ModelCallTrace, zdanie: string, udane: boolean) => void,
): OcenaWywolania {
  const trafnosc = wybor(
    'Trafność odpowiedzi modelu',
    SKALA.map((wartosc) => [
      wartosc,
      wartosc === ModelCallQuality.Unrated
        ? `${OCENA_WYWOLANIA[wartosc]} — zdejmij ocenę nadaną wcześniej`
        : OCENA_WYWOLANIA[wartosc],
    ]),
  );
  trafnosc.value = wywolanie.quality ?? ModelCallQuality.Unrated;

  const uzasadnienie = poleTresci(
    'Uzasadnienie oceny wywołania',
    2,
    'dlaczego odpowiedź jest taka; pole nieobowiązkowe',
  );
  uzasadnienie.value = wywolanie.qualityNote ?? '';

  const zapisz = przyciskAkcji('Zapisz ocenę wywołania', 'dn-btn dn-btn--atrament');
  zapisz.addEventListener('click', () => {
    const wybrana = trafnosc.value as ModelCallQuality;
    const nota = uzasadnienie.value.trim();
    void zrodlo
      .ocenWywolanie({
        callId: wywolanie.id,
        quality: wybrana,
        ...(nota === '' ? {} : { note: nota }),
      })
      .then((wynik) => {
        if (!wynik.udany || wynik.wynik === undefined) {
          const powod =
            wynik.blad === undefined
              ? ''
              : ` Powód: ${wynik.blad.message === '' ? 'rdzeń nie podał przyczyny' : wynik.blad.message} (kod ${wynik.blad.code}).`;
          poZapisie(
            wywolanie,
            `${zdanieNiepowodzenia('zapisu oceny wywołania', wynik.powod)}${powod}`,
            false,
          );
          return;
        }
        poZapisie(wynik.wynik, zdanieZapisu(wynik.wynik, wybrana, nota), true);
      });
  });

  const element = document.createElement('div');
  element.className = 'dg-narzedzie__pasek';
  element.dataset['ocenaWywolania'] = wywolanie.id;
  element.append(
    wiersz('Trafność odpowiedzi modelu', trafnosc, {
      klasa: 'dg-wiersz',
      objasnienie:
        'Ocena jest sądem Operatora zapisywanym w rdzeniu komendą provenance.call.rate, ' +
        'nie pomiarem. Skala pochodzi z kontraktu; okno jej nie wymyśla.',
    }),
    wiersz('Uzasadnienie oceny', uzasadnienie, {
      klasa: 'dg-wiersz',
      objasnienie: 'Pole nieobowiązkowe — kontrakt przyjmuje ocenę bez uzasadnienia.',
    }),
    zapisz,
  );

  return { element };
}

/**
 * Zdanie o skutku zapisu — porównanie zamówienia z tym, co rdzeń oddał.
 *
 * Rdzeń oddaje wywołanie po zapisie, więc podmiana oceny albo pominięcie
 * uzasadnienia są tu widoczne. Samo „ocena zapisana" byłoby zdaniem, po którym
 * Operator nadal nie wie, co w rdzeniu stoi.
 */
function zdanieZapisu(
  oddane: ModelCallTrace,
  zamowiona: ModelCallQuality,
  zamowionaNota: string,
): string {
  const oddanaOcena = oddane.quality ?? ModelCallQuality.Unrated;
  if (oddanaOcena !== zamowiona) {
    return (
      `Rdzeń zapisał ocenę INNĄ niż wysłana: wysłano „${OCENA_WYWOLANIA[zamowiona]}", ` +
      `oddał „${OCENA_WYWOLANIA[oddanaOcena]}".`
    );
  }
  const oddanaNota = oddane.qualityNote ?? '';
  if (zamowionaNota !== '' && oddanaNota !== zamowionaNota) {
    return (
      `Rdzeń zapisał ocenę „${OCENA_WYWOLANIA[oddanaOcena]}", ale uzasadnienia NIE zapisał ` +
      'w wysłanej postaci — sprawdź je po ponownym odczycie wywołania.'
    );
  }
  return (
    `Rdzeń zapisał ocenę wywołania: ${OCENA_WYWOLANIA[oddanaOcena]}` +
    (oddanaNota === '' ? ' bez uzasadnienia.' : ` wraz z uzasadnieniem (${String(oddanaNota.length)} znaków).`)
  );
}
