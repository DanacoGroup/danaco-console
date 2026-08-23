import { MODUL as MODUL_AGENTS } from '../moduly/agents/indeks';
import { testowanyAgent } from '../moduly/agents/testowany-agent';
import { politykaModulu, type PolitykaUlotnosci } from '../rozmowa/ulotnosc';

/**
 * Złożenie polityki pamięci rozmowy dla modułu okna.
 *
 * To, czy moduł ma pamięć sesyjną, niesie profil modułu
 * (`okno-komunikacji/profil-modulu.ts`, pole `pamiecSesyjna`) i wyłącznie
 * stamtąd bierze tę informację `politykaModulu`. Tutaj dochodzi jedna rzecz,
 * której warstwa rozmowy znać nie może: co w danym module jest kontekstem
 * roboczym, którego zmiana kończy rozmowę.
 *
 * Taki kontekst ma jeden moduł — Agents, gdzie kontekstem jest testowany
 * ekspert. Powłoka jest jedynym miejscem, w którym obie strony są widoczne
 * naraz: widok modułu powstaje w `przestrzen-modulu`, a rozmowa okna
 * w `wiazanie-gniazda`, i żadna z tych warstw nie zna drugiej.
 *
 * Moduł nieobjęty tą funkcją też przechodzi przez `politykaModulu` — dostaje
 * politykę wprost z własnego profilu, bez kontekstu roboczego. Odjęcie pamięci
 * sesyjnej kolejnemu modułowi czyni jego okno ulotnym bez zmiany w tym pliku.
 */
export function politykaOknaModulu(kod: string): PolitykaUlotnosci {
  return politykaModulu(kod, kod === MODUL_AGENTS.kod ? testowanyAgent : null);
}
