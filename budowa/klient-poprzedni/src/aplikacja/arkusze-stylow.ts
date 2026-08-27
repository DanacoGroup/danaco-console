/**
 * Kolejność wczytania warstw wizualnych interfejsu.
 *
 * Importy tworzą łańcuch zależności: żetony motywu dają barwy, typografię
 * i przestrzeń, biblioteka `.dn-*` stoi na nich, a `aplikacja.css` osadza
 * reguły w obszarze roboczym powłoki i idzie ostatni.
 */

import '../motyw/motyw.css';
import '../komponenty/indeks.css';
import '../okno-komunikacji/okno.css';
import './powloka.css';
import './aplikacja.css';
