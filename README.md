![Go](https://github.com/fogonthedowns/diplomatic/workflows/Go/badge.svg)

0) Loop through moves and determine if the destination is contested

1) ✅ Resolve all uncontested moves
   a) ✅ is the move possible? If so Update the move.ResolvedMove
   b) ✅ if not? Hold 

2) ✅ Collect all support orders
   a) ✅ determine if the support has been cut.

3) ✅ Collect all convoy orders
   a) ✅ determine if the path has been disrupted
   b) ✅ determine if the path is complete

4) resolve all contested moves
   a) ✅ will the piece move out of the way?
   b) ✅ is there support?

5) Edge cases
   a) ✅ Spain
   b) ✅ Bulgaria
   c) ✅ St. Petersburg

6) Retreat phase:
  a) ✅ Add dislodged concept to piece.
  b) ✅ If a piece is dislodged and it is retreat phase, it must move or be destroyed. 
  c) ✅ A unit can not retreat from the origin of the attack.
  d) ✅ Non dislodged pieces may not move in retreat phase.
  e) ✅ Update the territory.Owner by move.PieceOwner after fall retreat phase is adjudicated
  f) ✅ make Phase a type rather than an int (see newPhase, phase and Phase)
  g) when assigning edge case Territory Save, ensure master territory is assigned ownership

7) Build Phase:
  a) ✅ Add concept of Victory Center, count points
  b) ✅ Count active units
     i) if active units exceeds Victory Centers a destroy must be issued
    ii) if active units is less than Victory Center count a build order can be issued - only on home centers
  c) If a Destroy must be issued and one is not issued, destroy a random unit
  d) Create a new unit on the map, if a successful build order is issued
  e) Ensure the year is cycled at the end of the build phase

8) Victory:
  a) if the centers exceeds 18 at the end of the year the player wins the game
  b) a victory can be shared between players who agree to share a victory

9) Validations:
    a) ✅ verify the piece exists
    b) ✅ Verify you own the piece.
    c) ✅ verify the piece is at the matching start location 

design considerations:

✅ you can not attack your own territory
✅ you must be able to issue support and convoy orders for others
✅ Players who have not issued moves, will issue hold orders
Move accepts phase, should it? Should it be a lookup on game?
