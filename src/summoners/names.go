package summoners

import "strings"


const namesVariationWidth = 5

var ChampionNames = strings.Split(champions, "\n")

var NamesRange = len(ChampionNames) / namesVariationWidth

var champions = `Ahri
Akali
Alistar
Amumu
Anivia
Annie
Ashe
Blitzcrank
Brand
Caitlyn
Cassiopeia
Cho'Gath
Corki
Dr. Mundo
Evelynn
Ezreal
Fiddlesticks
Fizz
Galio
Gangplank
Garen
Gragas
Graves
Heimerdinger
Irelia
Janna
Jarvan IV
Jax
Karma
Karthus
Kassadin
Katarina
Kayle
Kennen
Kog'Maw
LeBlanc
Lee Sin
Leona
Lux
Malphite
Malzahar
Maokai
Master Yi
Miss Fortune
Mordekaiser
Morgana
Nasus
Nidalee
Nocturne
Nunu
Olaf
Orianna
Pantheon
Poppy
Rammus
Renekton
Riven
Rumble
Ryze
Sejuani
Shaco
Shen
Shyvana
Singed
Sion
Sivir
Skarner
Sona
Soraka
Swain
Talon
Taric
Teemo
Tristana
Trundle
Tryndamere
Twisted Fate
Twitch
Udyr
Urgot
Vayne
Veigar
Viktor
Vladimir
Volibear
Warwick
Wukong
Xerath
Xin Zhao
Yorick
Zilean`
