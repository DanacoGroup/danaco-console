import { elementIkony, type NazwaIkony } from '../ikony/ikony';

/**
 * Kafel jednej liczby telemetrycznej — para „liczba + co ona mierzy".
 *
 * Liczba idzie krojem technicznym (`--dn-ff-mono`), bo jest daną techniczną,
 * a nie nagłówkiem; krój szeryfowy zostaje przy tytułach sekcji. Etykieta
 * nazywa miarę wprost i nie bywa zastępowana samym kolorem.
 */
export interface OpisKafla {
  /** Wartość widoczna w kaflu. */
  wartosc: string;
  /** Co dokładnie mierzy liczba. */
  etykieta: string;
  /** Zdanie doprecyzowujące, widoczne pod etykietą i w podpowiedzi. */
  objasnienie?: string;
  /** Ikona wiodąca kafla. */
  ikona?: NazwaIkony;
  /** Czy kafel ma nieść akcent — zarezerwowane dla miary krytycznej. */
  wyrozniony?: boolean;
}

/** Buduje kafel liczby telemetrycznej. */
export function utworzKafelLiczby(opis: OpisKafla): HTMLElement {
  const kafel = document.createElement('div');
  kafel.className = 'mc-kafel';
  if (opis.wyrozniony) {
    kafel.classList.add('mc-kafel--wyrozniony');
  }
  if (opis.objasnienie) {
    kafel.title = opis.objasnienie;
  }

  const gora = document.createElement('div');
  gora.className = 'mc-kafel__gora';

  if (opis.ikona) {
    gora.append(elementIkony(opis.ikona, { rozmiar: 18, klasa: 'dn-ikona mc-kafel__ikona' }));
  }

  const wartosc = document.createElement('span');
  wartosc.className = 'mc-kafel__wartosc';
  wartosc.textContent = opis.wartosc;
  gora.append(wartosc);

  const etykieta = document.createElement('span');
  etykieta.className = 'mc-kafel__etykieta';
  etykieta.textContent = opis.etykieta;

  kafel.append(gora, etykieta);

  if (opis.objasnienie) {
    const objasnienie = document.createElement('span');
    objasnienie.className = 'mc-kafel__objasnienie';
    objasnienie.textContent = opis.objasnienie;
    kafel.append(objasnienie);
  }

  return kafel;
}
