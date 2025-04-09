import React from "react";
import { createIcon } from "@/components/ui/icon";
import { Path } from "react-native-svg";

function RivallBlackCrownFunc() {
  const icon = createIcon({
    viewBox: "0 0 2000 2000",
    path: (
      <>
        <Path d="" stroke="none" fill="#080404" fill-rule="evenodd"/>
        <Path d="M 985.500 594.648 C 962.656 596.431, 945.493 600.049, 928.857 606.588 C 894.333 620.161, 864.041 644.297, 844.254 674 C 816.337 715.908, 807.148 766.555, 819.095 812.675 C 823 827.755, 826.212 836.512, 832.437 849.056 C 845.495 875.369, 865.004 898.713, 888.730 916.414 C 893.807 920.202, 898.201 924.058, 898.495 924.984 C 899.242 927.337, 896.858 933.607, 891.346 943.792 C 866.707 989.314, 834.099 1026.722, 794 1055.468 C 752.983 1084.872, 711.740 1102.415, 657.500 1113.531 C 638.924 1117.338, 598.040 1118.938, 570 1116.955 C 552.251 1115.700, 546.876 1114.888, 526.500 1110.385 C 478.224 1099.717, 430.617 1076.779, 391.500 1045.341 C 378.959 1035.261, 355.979 1012.849, 347.632 1002.555 C 340.459 993.711, 340.409 992.033, 346.946 979.775 C 361.300 952.857, 366.961 926.016, 365.667 891 C 364.911 870.546, 363.635 862.218, 358.895 846.780 C 349.784 817.107, 335.815 793.336, 315.106 772.268 C 290.203 746.932, 262.149 731.676, 224 722.725 C 210.762 719.619, 210.015 719.557, 185.500 719.529 C 164.752 719.505, 159.060 719.816, 152.033 721.360 C 95.046 733.880, 47.783 770.602, 24.437 820.500 C 10.413 850.474, 4.186 882.830, 6.796 912.171 C 11.961 970.234, 42.955 1021.474, 92.252 1053.452 C 115.004 1068.210, 145.293 1078.126, 175.329 1080.649 C 188.617 1081.765, 187.227 1070.950, 186.944 1171.036 C 186.260 1413.220, 185.597 1389.655, 193.316 1397.363 C 198.868 1402.907, 203.886 1404.354, 214.924 1403.595 C 223.935 1402.975, 404.743 1402.234, 719 1401.528 C 822.125 1401.296, 989.300 1400.855, 1090.500 1400.548 C 1191.700 1400.241, 1392.400 1399.947, 1536.500 1399.893 C 1782.628 1399.803, 1798.742 1399.692, 1802.500 1398.065 C 1804.700 1397.113, 1807.943 1394.690, 1809.707 1392.681 C 1816.092 1385.409, 1815.825 1394.428, 1814.505 1230.500 C 1813.845 1148.550, 1813.497 1080.969, 1813.732 1080.320 C 1814.301 1078.753, 1819.495 1077, 1823.570 1077 C 1825.369 1077, 1833.499 1075.872, 1841.637 1074.494 C 1874.318 1068.957, 1901.820 1056.262, 1927.500 1034.858 C 1971.735 997.989, 1996.527 940.401, 1992.950 882.832 C 1990.179 838.246, 1967.224 790.629, 1934.589 761.771 C 1912.976 742.660, 1892.636 731.022, 1866.238 722.663 C 1845.486 716.092, 1841.119 715.496, 1814 715.538 C 1789.896 715.575, 1789.282 715.628, 1776 718.809 C 1757.675 723.199, 1748.561 726.281, 1735.031 732.661 C 1686.245 755.667, 1652.006 798.143, 1637.986 853.053 C 1635.255 863.752, 1634.035 876.704, 1634.035 895 C 1634.035 928.271, 1638.142 946.310, 1652.117 974.408 C 1656.921 984.067, 1658.841 988.956, 1658.403 990.408 C 1655.489 1000.056, 1618.023 1036.756, 1592.002 1055.450 C 1572.472 1069.481, 1542.820 1085.468, 1519.903 1094.321 C 1496.365 1103.415, 1464.208 1111.604, 1441.803 1114.209 C 1428.591 1115.745, 1371.744 1115.753, 1358.500 1114.221 C 1343.239 1112.456, 1316.262 1106.364, 1299.106 1100.810 C 1268.023 1090.745, 1242.892 1078.709, 1216.698 1061.340 C 1195.402 1047.219, 1183.253 1037.331, 1165.954 1020.040 C 1144.025 998.121, 1128.576 977.913, 1112.790 950.500 C 1103.475 934.324, 1100.175 926.543, 1101.633 924.192 C 1102.210 923.261, 1106.466 919.582, 1111.091 916.015 C 1147.778 887.721, 1170.867 851.852, 1181.661 806.383 C 1184.355 795.033, 1184.478 793.445, 1184.405 771 L 1184.328 747.500 1180.304 732 C 1175.659 714.112, 1172.056 704.616, 1164.443 690.197 C 1148.803 660.579, 1127.244 637.966, 1097.538 620.023 C 1080.410 609.677, 1056.311 600.396, 1039.023 597.488 C 1022.359 594.684, 1000.146 593.506, 985.500 594.648" stroke="none" fill="#040404" fill-rule="evenodd"/>
      </>
    )
  })
  return icon
}

const RivallBlackCrown = RivallBlackCrownFunc()
export { RivallBlackCrown }