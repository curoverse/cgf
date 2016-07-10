// query.js
//

print = ((typeof(print)==="undefined") ? console.log : print);

var cgf_info = {};

function setup_cgf_info() {
  cgf_info.cgf = [];
  cgf_info.cgf.push({ "file":"../data/hu826751-GS03052-DNA_B01.cgf", "name":"hu826751-GS03052-DNA_B01", "id":0 });
  cgf_info.cgf.push({ "file":"../data/hu0211D6-GS01175-DNA_E02.cgf", "name":"hu0211D6-GS01175-DNA_E02", "id":1 });

  cgf_info.id = {};
  for (var idx=0; idx<cgf_info.cgf.length; idx++) {
    cgf_info.id[ cgf_info.cgf[idx].name ] = idx;
  }

  cgf_info["CGFVersion"] = "0.1.0";
  cgf_info["CGFLibVersion"] = "0.1.0";
  cgf_info["PathCount"] = 863;
  cgf_info["StepPerPath"] = [
     5433, 11585, 7112, 7550, 13094, 10061, 15111, 13212, 14838, 7361, 8565, 8238, 21058, 15318, 9982, 14543,
    20484, 11704, 9056, 29572, 3032, 58941, 13626, 13753, 10082, 19756, 9669, 18011, 17221, 16418, 6572, 10450,
    653, 1, 1, 43, 4603, 4524, 17225, 5245, 9951, 5416, 18877, 6467, 14301, 7627, 11539, 16593,
    21475, 19845, 11886, 19126, 30932, 16774, 11607, 37511, 1368, 9016, 14132, 15803, 6847, 26570, 19594, 17082,
    10529, 20354, 17716, 9931, 19189, 14703, 8418, 8231, 17045, 7804, 12459, 23570, 20025, 8246, 24611, 10263,
    17693, 11001, 7904, 5629, 32719, 19083, 565, 3431, 20757, 13319, 5383, 9608, 10026, 16921, 14381, 29377,
    6845, 8754, 6367, 21554, 7707, 18707, 4227, 2345, 16932, 19091, 15332, 23909, 32173, 10128, 9612, 24819,
    9782, 21619, 22599, 5851, 16177, 24645, 24453, 14657, 3551, 19209, 17178, 6784, 22677, 10729, 4764, 18388,
    11981, 5804, 12040, 29022, 9918, 17574, 4842, 16740, 11327, 16335, 1542, 416, 23880, 6126, 8255, 16187,
    20267, 23705, 17658, 21050, 14728, 14705, 2708, 9599, 1327, 17097, 6536, 3446, 7194, 13517, 6740, 12960,
    8454, 15276, 6666, 10736, 7497, 7113, 13394, 16658, 7897, 10893, 15843, 24193, 12589, 10989, 7735, 7704,
    6591, 26835, 12945, 19129, 12707, 14282, 6739, 5660, 7363, 17599, 20166, 15899, 5832, 18674, 15349, 10225,
    13863, 25249, 32580, 20511, 13259, 14135, 3468, 71, 25343, 27513, 14097, 21456, 9860, 13680, 6387, 10838,
    4120, 21815, 5451, 14460, 8533, 24975, 24610, 25300, 11590, 19404, 8688, 32414, 7729, 19437, 6621, 10118,
    17649, 24182, 10736, 21411, 6710, 17505, 4790, 22874, 15243, 14561, 17381, 7292, 13961, 20750, 12771, 17639,
    5133, 16978, 18906, 16519, 15821, 13209, 882, 4225, 31741, 15233, 1182, 13597, 6528, 11710, 13632, 16991,
    5455, 37078, 22890, 16898, 6764, 20266, 7277, 6180, 8009, 24144, 22877, 12483, 21662, 12287, 19473, 20872,
    11085, 11566, 16415, 34070, 16922, 13794, 14120, 8663, 7451, 11295, 13748, 3815, 7213, 7030, 38651, 6143,
    12781, 5883, 5178, 11753, 15562, 22214, 22047, 4132, 16117, 3941, 144, 4865, 400, 25489, 22288, 30139,
    3706, 11083, 19909, 24752, 4171, 19061, 35002, 14079, 794, 29730, 3892, 12776, 3515, 15587, 14919, 14827,
    11010, 13427, 13368, 11662, 21111, 13834, 24662, 10333, 6684, 8376, 25611, 10830, 17440, 17699, 9856, 3300,
    23551, 7908, 24000, 7739, 13746, 5876, 13653, 11619, 105, 227, 12689, 19087, 11490, 35461, 6928, 11137,
    6317, 19717, 18677, 2636, 10982, 28108, 11243, 14787, 10618, 12904, 7678, 4053, 8783, 21899, 18003, 16798,
    16058, 8543, 15728, 7511, 16071, 18591, 25102, 17085, 16227, 5457, 29901, 6958, 5306, 12761, 2290, 4222,
    15593, 1523, 10990, 23625, 2365, 14954, 7597, 9733, 12983, 17099, 7155, 17446, 7771, 24670, 22012, 9790,
    17944, 16958, 6352, 22341, 6025, 12803, 18803, 16509, 19724, 13970, 23963, 7842, 9501, 16725, 20807, 9222,
    7462, 5182, 22155, 9365, 20144, 11012, 8142, 1490, 180, 546, 1, 1, 15, 550, 4865, 7015,
    20266, 7250, 11850, 10403, 13346, 5036, 7311, 10212, 9994, 12206, 21611, 12006, 13925, 10860, 19459, 12846,
    17584, 11203, 1904, 7356, 5714, 14022, 11522, 3238, 10867, 22206, 19356, 3286, 381, 14758, 7681, 18901,
    6319, 11569, 13319, 2602, 1, 12601, 5388, 8544, 32551, 13246, 23124, 16676, 10420, 16083, 23002, 4756,
    13393, 4473, 10500, 8904, 9750, 4253, 7078, 3459, 24069, 12012, 16737, 10252, 5577, 17329, 11901, 19092,
    9991, 28650, 8063, 13688, 21339, 17049, 4291, 15046, 21055, 27571, 19581, 5339, 1, 2796, 15653, 6733,
    5702, 9463, 8431, 7485, 17429, 7445, 33236, 10017, 15088, 16390, 18985, 3047, 29163, 8290, 8000, 26700,
    10459, 15540, 11802, 16858, 12184, 8407, 15777, 9945, 7774, 20407, 5030, 20355, 4994, 11256, 9088, 5210,
    703, 31263, 9981, 8655, 12869, 6059, 5323, 19308, 6962, 10252, 14659, 16466, 18159, 25083, 8822, 14458,
    13654, 20804, 8472, 20356, 9936, 2048, 7595, 10099, 4973, 9834, 18782, 13534, 16861, 1, 1, 1,
    1, 449, 13648, 8140, 8894, 4307, 12796, 7164, 5979, 18211, 19843, 2279, 5677, 13654, 16553, 17021,
    10676, 13343, 11629, 19081, 8331, 7079, 7216, 33870, 9290, 20014, 12554, 4179, 9303, 12659, 8980, 13317,
    17551, 1, 1, 1, 1, 91, 16293, 33478, 7694, 4755, 4736, 21768, 13932, 14148, 12245, 5458,
    10017, 15321, 10317, 11761, 9101, 13816, 21162, 17182, 5312, 19338, 8096, 10791, 6468, 20877, 6861, 3000,
    11596, 1, 1, 1, 1, 650, 9508, 9670, 6240, 919, 7453, 25276, 10122, 2914, 3833, 18009,
    12803, 23978, 730, 17531, 13228, 383, 818, 20600, 9375, 4772, 6376, 13251, 9675, 14940, 19964, 17117,
    15125, 29145, 9758, 7794, 8325, 3452, 13738, 9153, 15025, 13310, 2344, 1, 1932, 22098, 16158, 2693,
    36810, 14074, 7919, 4845, 19451, 10051, 10058, 11572, 5454, 5493, 11041, 11843, 15854, 19846, 17827, 125,
    1653, 21451, 21850, 1084, 9274, 12463, 8800, 10895, 28728, 2071, 9705, 5530, 5548, 10683, 15160, 14696,
    1860, 22145, 10747, 16523, 5517, 9195, 13344, 12, 1771, 23326, 30739, 18023, 25450, 18584, 21768, 9509,
    10948, 10287, 21091, 7440, 17747, 19563, 23601, 23077, 347, 7918, 12998, 13442, 559, 2780, 15135, 11458,
    9677, 1431, 16004, 5610, 9780, 11468, 7764, 8969, 10185, 19284, 16238, 11893, 23036, 13336, 3819, 12729,
    1268, 1, 8908, 8382, 11966, 16626, 1577, 16554, 12901, 20235, 6003, 7836, 17926, 1, 1, 1214,
    534, 1, 4755, 30050, 11265, 18230, 16716, 7896, 7554, 11599, 21326, 1, 1, 1, 1, 2745,
    11905, 4602, 8007, 14401, 9459, 20828, 12737, 11643, 16587, 4104, 6858, 5235, 6576, 13137, 29283, 8472,
    9959, 11291, 16995, 8588, 23499, 18569, 14851, 10837, 14462, 10224, 1492, 3714, 5149, 10944, 13980, 6118,
    6368, 29977, 5799, 12454, 4748, 18033, 14477, 3916, 18518, 28427, 15228, 29028, 6516, 11944, 15846, 8098,
    6040, 18525, 26363, 1, 1045, 12490, 1, 361, 4758, 14711, 4019, 5647, 721, 181, 35
  ];

  return cgf_info;
}
setup_cgf_info();

function tile_concordance_slice(set_a, set_b, lvl) {
}

function tile_concordance(set_a, set_b, lvl) {
  var A = [], B = [];

  lvl=((typeof(lvl)==="undefined")?2:lvl);
  if (lvl<0) { lvl=0; }
  if (lvl>2) { lvl=2; }

  for (var i=0; i<set_a.length; i++) {
    if (typeof(set_a[i])==="string") {
      if (set_a[i] in cgf_info.id) {
        A.push(cgf_info.id[set_a[i]]);
      } else if ((set_a[i]>=0) && (set_a[i]<cgf_info.id.length)) {
        A.push(set_a[i]);
      }
    }
  }

  for (var i=0; i<set_b.length; i++) {
    if (typeof(set_b[i])==="string") {
      if (set_b[i] in cgf_info.id) {
        B.push(cgf_info.id[set_b[i]]);
      } else if ((set_b[i]>=0) && (set_b[i]<cgf_info.id.length)) {
        B.push(set_b[i]);
      }
    }
  }

  var rstr = "";
  //var rstr = muduk_cgf_tile_concordance(A,B,lvl);

  var r = {};
  try {
    r = JSON.parse(rstr);
  } catch(err) {
    r["result"] = "error: parse error" + String(err);
  }

  return r;
}

function help() {
  print("muduk server");
}

function query(q) {
  var qobj = {};
  var robj = {};
  try {
    qobj = JSON.parse(q);
  } catch(err) {
    return err;
  }

  if ("request" in qobj) {
    print("request: " + String(qobj["request"]));
  }

  return JSON.stringify(robj);
}

function muduk_return(q, indent) {
  indent = ((typeof(indent)==="undefined") ? '' : indent);
  if (typeof(q)==="undefined") { return ""; }
  if (typeof(q)==="object") {
    var s = "";
    try {
      s = JSON.stringify(q, null, indent);
    } catch(err) {
    }
    return s;
  }
  if (typeof(q)==="string") { return q; }
  if (typeof(q)==="number") { return q; }
  return "";
}
