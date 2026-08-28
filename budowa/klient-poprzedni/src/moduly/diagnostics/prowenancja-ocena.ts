import { ModelCallQuality, type ModelCallTrace } from '../../../../shared/contract';
import { poleTresci, przyciskAkcji, wiersz, wybor } from '../../modele/kontrolki-formularza-braki';
import { zdanieNiepowodzenia } from './niepowodzenie-odczytu';
import { OCENA_WYWOLANIA } from './prowenancja-slowa';
import type { ZrodloProwenancji } from './prowenancja-zrodlo';

/**
 * Ocena wywołania modelu jest czynnością zapisu, a nie odczytem: komenda
 * `provenance.call.rate` zapisuje sąd człowieka o pracy modelu. Okno daje jej
 * jawny chwyt złożony z wyboru trafności, pola uzasadnienia i przycisku zapisu.
 */
export interface OcenaWywolania {
  element: HTMLElement;
}

/**
 * Trafności w kolejności od najlepszej, wzięte z wyliczenia ModelCallQuality
 * kontraktu; wartość „bez oceny” stoi na końcu, ponieważ służy zdjęciu oceny
 * nadanej wcześniej.
 */
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
 * Zdanie o skutku zapisu oceny: zestawia zamówienie z tym, co rdzeń oddał.
 * Rdzeń oddaje wywołanie po zapisie, więc podmiana oceny albo pominięcie
 * uzasadnienia są w tym zdaniu widoczne.
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
